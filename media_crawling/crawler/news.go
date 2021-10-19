package crawler

import (
	"encoding/json"
	"fmt"
	"media_crawling/config"
	"media_crawling/models"
	"media_crawling/util"
	"net/url"
	"regexp"
	"strings"
)

// requestNaverNews : 네이버 뉴스 API에 요청
func requestNaverNews(query string) NaverNewsResponse {
	encodedQuery := url.QueryEscape(query)
	requestURL := "https://openapi.naver.com/v1/search/news.json"
	completedURL := fmt.Sprintf("%s?query=%s&display=100&start=1&sort=date", requestURL, encodedQuery)

	var headers map[string]string
	headers = make(map[string]string)
	headers["X-Naver-Client-Id"] = config.Secret["naver:id"]
	headers["X-Naver-Client-Secret"] = config.Secret["naver:secret"]
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	result := Request(completedURL, headers)
	news := NaverNewsResponse{}
	json.Unmarshal([]byte(result), &news)

	return news
}

// refineWeirdLink : 이상한 링크 정제
// 영남일보는 http://나 https://가 아닌 ://만 붙여서 나온다.
// 이와 같은 예외 상황을 일부 처리한다.
func refineWeirdLink(link string) string {
	// key는 이상한 link, value는 수정한 link
	weirdCaeses := map[string]string{
		"//www.yeongnam.com":    "https://www.yeongnam.com",
		"http:/www.segyefn.com": "http://www.segyefn.com",
		"goal.com/kr":           "https://www.goal.com/kr",
	}

	// http://나 https://로 시작하는지 확인
	// 맞으면 입력 받은 link 그대로 리턴, 아니면 수정한 link 리턴
	pattern := `^(?:http://|https://)`
	re := regexp.MustCompile(pattern)
	match := string(re.Find([]byte(link)))
	if match != "" {
		return link
	}

	// 수정
	for weirdCase, correctedCase := range weirdCaeses {
		link = strings.Replace(link, weirdCase, correctedCase, 1)
	}

	return link
}

// getPressKeyword : 언론사 키워드
// pattern 매칭에 실패하면 ""(empty string)을 리턴
func getPressKeyword(link string) string {
	pattern := `^(?:http://|https://)([^/]*)`
	re := regexp.MustCompile(pattern)
	refinedLink := refineWeirdLink(link)
	match := string(re.Find([]byte(refinedLink)))

	if match != "" {
		match = strings.ReplaceAll(match, "http://", "")
		match = strings.ReplaceAll(match, "https://", "")
	}

	return match
}

// refineNews : 네이버 뉴스 API 응답 값을 정제
func refineNews(presses []models.Press, news NaverNewsResponse, query string) []models.News {
	var refinedNewsList []models.News

	// uselessWords에서 key는 uselessWord, value는 교체값이다.
	uselessWords := GetUselessWords()

	for _, item := range news.Items {
		item.Title = util.ReplaceString(item.Title, uselessWords)
		item.Description = util.ReplaceString(item.Description, uselessWords)
		PressKeyword := getPressKeyword(item.Originallink)
		// item.PubDate는 timezone이 표기되는데, 기본 형변환 시 timezone까지 고려되지 않아 함수를 이용하여 따로 처리
		datetime := util.GetKoreanDateTime(item.PubDate, "RFC1123Z", false, DatetimeFormat)
		refinedNewsList = append(refinedNewsList, models.News{
			Query:        query,
			Title:        item.Title,
			Link:         item.Originallink,
			PressKeyword: PressKeyword,
			Description:  item.Description,
			PubDate:      datetime})
	}

	return refinedNewsList
}

// CrawlNews : naver news API를 이용하여 뉴스 기사를 크롤링한 후 저장합니다.
func CrawlNews(logFilePath string) {
	fpLog, logger := util.GetFileLogger(logFilePath)
	defer util.CloseFileLogger(fpLog)

	db := GetDB("operation")
	defer CloseDB(db)

	presses := models.GetPresses(db)
	newsQueries := models.GetQueries(db, "news", 9999)

	logger.Println("뉴스기사 수집 정상시작")
	for _, newsQuery := range newsQueries {
		logger.Printf("검색어 \"%s\"에 관한 수집 시작\n", newsQuery)
		news := requestNaverNews(newsQuery)

		// 정제
		refinedNewsList := refineNews(presses, news, newsQuery)

		// refinedNewsList를 ImportDataToDB의 data 파라미터로 사용하기 위해
		// data 변수를 선언하고, refinedNewsList의 각 행을 할당한다.
		data := make([]interface{}, len(refinedNewsList))
		for index, row := range refinedNewsList {
			data[index] = row
		}
		// DB에 업로드
		ImportDataToDB(db, "news", "", data)
		logger.Printf("검색어 \"%s\"에 관한 수집 종료\n", newsQuery)
	}

	logger.Println("뉴스기사 수집 정상종료")
}
