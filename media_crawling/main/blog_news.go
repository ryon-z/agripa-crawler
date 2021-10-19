package main

import (
	"media_crawling/checker"
	"media_crawling/config"
	"media_crawling/crawler"
	"media_crawling/util"
	"os"
)

func main() {
	conf := config.Conf
	newsLogFilePath := util.MakePath([]string{conf.SubLogDirPath, "news.txt"})
	newsLogErrorMessage := "뉴스 수집 문제 생김"

	for _, path := range []string{newsLogFilePath} {
		os.Remove(path)
	}
	util.MakeDirIfNotExists(conf.LogDirPath)
	util.MakeDirIfNotExists(conf.SubLogDirPath)

	crawler.CrawlNews(newsLogFilePath)
	checker.CheckLog(newsLogFilePath, newsLogErrorMessage)

	//// 블로그글도 불필요한 자료가 많아 다운로드 보류
	// blogLogFilePath := util.MakePath([]string{conf.SubLogDirPath, "blog.txt"})
	// blogLogErrorMessage := "블로그 수집 문제 생김"
	// crawler.CrawlBlogs(blogLogFilePath)
	// checker.CheckLog(blogLogFilePath, blogLogErrorMessage)

	//// 카페글은 게시일이 표기되지 않아 다운로드 보류
	// cafeLogFilePath := util.MakePath([]string{config.Conf.SubLogDirPath, "cafe.txt"})
	// cafeLogErrorMessage := "카페 수집 문제 생김"
	// crawler.CrawlCafes(cafeLogFilePath, cafeLogErrorMessage)
}
