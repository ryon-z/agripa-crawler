package crawler

import (
	"encoding/json"
	"fmt"
	"media_crawling/config"
	"media_crawling/models"
	"media_crawling/util"
	"net/url"
)

func requestNaverCafes(query string) NaverCafeResponse {
	encodedQuery := url.QueryEscape(query)
	requestURL := "https://openapi.naver.com/v1/search/cafearticle.json"
	completedURL := fmt.Sprintf("%s?query=%s&display=100&start=1&sort=date", requestURL, encodedQuery)

	var headers map[string]string
	headers = make(map[string]string)
	headers["X-Naver-Client-Id"] = config.Secret["naver:id"]
	headers["X-Naver-Client-Secret"] = config.Secret["naver:secret"]
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	result := Request(completedURL, headers)
	cafes := NaverCafeResponse{}
	json.Unmarshal([]byte(result), &cafes)

	return cafes
}

func refineCafes(cafes NaverCafeResponse, query string) []models.Cafe {
	var refinedCafes []models.Cafe

	// uselessWords에서 key는 uselessWord, value는 교체값이다.
	uselessWords := GetUselessWords()

	for _, item := range cafes.Items {
		item.Title = util.ReplaceString(item.Title, uselessWords)
		item.Description = util.ReplaceString(item.Description, uselessWords)
		item.Cafename = util.ReplaceString(item.Cafename, uselessWords)
		refinedCafes = append(refinedCafes, models.Cafe{
			Query:       query,
			Title:       item.Title,
			Link:        item.Link,
			Cafeurl:     item.Cafeurl,
			Description: item.Description,
			Cafename:    item.Cafename})
	}

	return refinedCafes
}

// CrawlCafes : 네이버 카페글을 수집합니다.
func CrawlCafes(logFilePath string) {
	fpLog, logger := util.GetFileLogger(logFilePath)
	defer util.CloseFileLogger(fpLog)

	db := GetDB("operation")
	defer CloseDB(db)

	cafeQueries := models.GetQueries(db, "cafe", 9999)

	logger.Println("카페글 수집 정상시작")
	for _, cafeQuery := range cafeQueries {
		logger.Printf("검색어 \"%s\"에 대한 수집 시작\n", cafeQuery)
		cafes := requestNaverCafes(cafeQuery)

		// 정제
		refinedCafes := refineCafes(cafes, cafeQuery)

		// DB에 업로드
		data := make([]interface{}, len(refinedCafes))
		for index, row := range refinedCafes {
			data[index] = row
		}
		ImportDataToDB(db, "cafe", "", data)
		logger.Printf("검색어 \"%s\"에 대한 수집 종료\n", cafeQuery)
	}

	logger.Println("카페글 수집 정상종료")
}
