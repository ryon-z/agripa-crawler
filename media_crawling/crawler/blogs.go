package crawler

import (
	"encoding/json"
	"fmt"
	"media_crawling/config"
	"media_crawling/models"
	"media_crawling/util"
	"net/url"
)

func requestNaverBlogs(query string) NaverBlogResponse {
	encodedQuery := url.QueryEscape(query)
	requestURL := "https://openapi.naver.com/v1/search/blog.json"
	completedURL := fmt.Sprintf("%s?query=%s&display=100&start=1&sort=date", requestURL, encodedQuery)

	var headers map[string]string
	headers = make(map[string]string)
	headers["X-Naver-Client-Id"] = config.Secret["naver:id"]
	headers["X-Naver-Client-Secret"] = config.Secret["naver:secret"]
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	result := Request(completedURL, headers)
	blogs := NaverBlogResponse{}
	json.Unmarshal([]byte(result), &blogs)

	return blogs
}

func refineBlogs(blogs NaverBlogResponse, query string) []models.Blog {
	var refinedBlogs []models.Blog

	// uselessWords에서 key는 uselessWord, value는 교체값이다.
	uselessWords := GetUselessWords()

	for _, item := range blogs.Items {
		item.Title = util.ReplaceString(item.Title, uselessWords)
		item.Description = util.ReplaceString(item.Description, uselessWords)
		item.Bloggername = util.ReplaceString(item.Bloggername, uselessWords)
		datetime := fmt.Sprintf("%s-%s-%s 00:00:01", item.Postdate[:4], item.Postdate[4:6], item.Postdate[6:])
		refinedBlogs = append(refinedBlogs, models.Blog{
			Query:       query,
			Title:       item.Title,
			Link:        item.Link,
			Bloggerlink: item.Bloggerlink,
			Description: item.Description,
			Bloggername: item.Bloggername,
			Postdate:    datetime})
	}

	return refinedBlogs
}

// CrawlBlogs : 네이버 블로그 글을 수집합니다.
func CrawlBlogs(logFilePath string) {
	fpLog, logger := util.GetFileLogger(logFilePath)
	defer util.CloseFileLogger(fpLog)

	db := GetDB("operation")
	defer CloseDB(db)

	blogQueries := models.GetQueries(db, "blog", 9999)

	logger.Println("블로그 수집 정상시작")
	for _, blogQuery := range blogQueries {
		logger.Printf("검색어 \"%s\"에 대한 수집 시작\n", blogQuery)
		blogs := requestNaverBlogs(blogQuery)

		// 정제
		refinedBlogs := refineBlogs(blogs, blogQuery)

		// DB에 업로드
		data := make([]interface{}, len(refinedBlogs))
		for index, row := range refinedBlogs {
			data[index] = row
		}
		ImportDataToDB(db, "blog", "", data)
		logger.Printf("검색어 \"%s\"에 대한 수집 종료\n", blogQuery)
	}

	logger.Println("블로그 수집 정상종료")
}
