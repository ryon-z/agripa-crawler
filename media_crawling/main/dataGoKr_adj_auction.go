package main

import (
	"media_crawling/checker"
	"media_crawling/config"
	"media_crawling/crawler"
	"media_crawling/refiner"
	"media_crawling/util"
	"os"
	"time"
)

func main() {
	conf := config.Conf
	loc := conf.Timezone
	logFilePath := util.MakePath([]string{conf.SubLogDirPath, "adjusted_auction.txt"})
	logErrorMessage := "정산가격 수집 문제 생김"

	os.Remove(logFilePath)
	util.MakeDirIfNotExists(conf.LogDirPath)
	util.MakeDirIfNotExists(conf.SubLogDirPath)

	// 어제 날짜
	start := time.Now().In(loc).AddDate(0, 0, -1)
	end := start

	// 임의 시간 설정
	/*
		start := time.Date(2020, 8, 18, 0, 0, 0, 0, loc)
		end := time.Date(2020, 8, 18, 0, 0, 0, 0, loc)
	*/

	// 조사가격과 매핑되는 표준품종코드만 API 요청(DEPRECIATED)
	// crawler.CrawlAuction("adjusted", start, end)

	// 모든 표준품종코드 API 요청
	crawler.CrawlAllAuctions("adjusted", start, end, logFilePath)
	checker.CheckLog(logFilePath, logErrorMessage)

	// 모든 표준품종코드 API 요청 결과에서 조사가격이 존재하는 row만 조회테이블로 저장
	refiner.RefineAcution("adjusted", start, end)
}
