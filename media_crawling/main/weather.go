package main

import (
	"media_crawling/checker"
	"media_crawling/config"
	"media_crawling/crawler"
	"media_crawling/util"
	"os"
	"time"
)

func main() {
	conf := config.Conf
	loc := conf.Timezone
	logFilePath := util.MakePath([]string{conf.SubLogDirPath, "weather.txt"})
	logErrorMessage := "일별 날씨 수집 문제 생김"

	os.Remove(logFilePath)
	util.MakeDirIfNotExists(conf.LogDirPath)
	util.MakeDirIfNotExists(conf.SubLogDirPath)

	// 어제 날짜
	start := time.Now().In(loc).AddDate(0, 0, -1)
	end := time.Now().In(loc)

	// 임의 기간 설정
	/*
		start := time.Date(2000, 1, 1, 0, 0, 0, 0, loc)
		end := time.Date(2020, 8, 23, 0, 0, 0, 0, loc)
	*/

	// 모든 일별 날씨 API 요청
	crawler.CrawlAllWeather(start, end, logFilePath)
	checker.CheckLog(logFilePath, logErrorMessage)
}
