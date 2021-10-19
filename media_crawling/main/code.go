package main

import (
	"fmt"
	"media_crawling/checker"
	"media_crawling/config"
	"media_crawling/crawler"
	"media_crawling/models"
	"media_crawling/refiner"
	"media_crawling/util"
	"os"
	"time"
)

func main() {
	conf := config.Conf
	_ = conf

	// 로그 처리(처리가 굉장히 복잡하기 때문에 main에서 로그 처리)
	logFilePath := util.MakePath([]string{conf.SubLogDirPath, "code.txt"})
	logIdentifier := "코드 갱신"
	logErrorMessage := fmt.Sprintf("%s 문제 생김")
	os.Remove(logFilePath)
	util.MakeDirIfNotExists(conf.LogDirPath)
	util.MakeDirIfNotExists(conf.SubLogDirPath)

	fpLog, logger := util.GetFileLogger(logFilePath)
	defer util.CloseFileLogger(fpLog)

	logger.Println(fmt.Sprintf("%s 정상시작", logIdentifier))

	// 다운로드
	crawler.CrawlAllGarakCodes()
	time.Sleep(time.Second * 3)
	crawler.CrawlAllWholesaleMarketCodes()
	time.Sleep(time.Second * 3)
	crawler.CrawlAllWholesaleMarketCoCodes()
	time.Sleep(time.Second * 3)
	crawler.CrawlAllStdGradeCode()
	time.Sleep(time.Second * 3)
	crawler.CrawlAllStdUnitCode()
	time.Sleep(time.Second * 3)
	crawler.CrawlAllPlaceOriginCode()
	time.Sleep(time.Second * 3)
	crawler.CrawlAllStdSpeciesCode()
	time.Sleep(time.Second * 3)

	// 정제
	refiner.SetStdItemKeword()
	refiner.SetItemMapping()
	refiner.SetItemCode()

	// 수집 DB에서 운영 DB로 테이블 교체
	crawler.MoveDataUsingDump("collection", "operation", models.StdItemKeyword{}.TableName(), false)
	crawler.MoveDataUsingDump("collection", "operation", models.ItemMapping{}.TableName(), false)
	crawler.MoveDataUsingDump("collection", "operation", models.StdItemCode{}.TableName(), false)

	// 에러 처리
	logger.Println(fmt.Sprintf("%s 정상종료", logIdentifier))
	checker.CheckLog(logFilePath, logErrorMessage)
}
