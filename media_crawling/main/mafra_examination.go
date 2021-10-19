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
	loc := conf.Timezone

	// 로그 처리(처리가 굉장히 복잡하기 때문에 main에서 로그 처리)
	logFilePath := util.MakePath([]string{conf.SubLogDirPath, "mafra_examination.txt"})
	logIdentifier := "도매가격과 소매가격 수집 및 정제"
	logErrorMessage := fmt.Sprintf("%s 문제 생김", logIdentifier)
	os.Remove(logFilePath)
	util.MakeDirIfNotExists(conf.LogDirPath)
	util.MakeDirIfNotExists(conf.SubLogDirPath)

	fpLog, logger := util.GetFileLogger(logFilePath)
	defer util.CloseFileLogger(fpLog)

	logger.Println(fmt.Sprintf("%s 정상시작", logIdentifier))

	// 어제 날짜
	start := time.Now().In(loc).AddDate(0, 0, -1)
	end := time.Now().In(loc)

	// 임의 시간 설정
	// (임의 시간 사용 시 기간이 길다면 DB Editor에서 수동으로 실행 후 dump 명령으로 옮길 것)
	// start := time.Date(2020, 9, 26, 0, 0, 0, 0, loc)
	// end := time.Date(2020, 10, 20, 0, 0, 0, 0, loc)

	startDate := util.GetDateString(start, "-")
	endDate := util.GetDateString(end, "-")
	dateColumnName := "ExaminDate"

	crawler.CrawlMafraExamination(start, end)
	refiner.UpdateMafraWholePrice(start, end)
	sqlQuery := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE %s >= "%s"
		AND %s <= "%s"
		;`, models.MafraWholePrice{}.TableName(),
		dateColumnName, startDate, dateColumnName, endDate)
	crawler.MoveDataUsingGorm("collection", "operation", "mafraWholePrice", sqlQuery)

	refiner.UpdateMafraRetailPrice(start, end)
	sqlQuery = fmt.Sprintf(`
		SELECT * FROM %s
		WHERE %s >= "%s"
		AND %s <= "%s"
		;`, models.MafraRetailPrice{}.TableName(),
		dateColumnName, startDate, dateColumnName, endDate)
	crawler.MoveDataUsingGorm("collection", "operation", "mafraRetailPrice", sqlQuery)

	// 에러 처리
	logger.Println(fmt.Sprintf("%s 정상종료", logIdentifier))
	checker.CheckLog(logFilePath, logErrorMessage)

	// 덤프 명령어
	// crawler.MoveDataUsingDump("collection", "operation", models.MafraWholePrice{}.TableName(), false)
	// crawler.MoveDataUsingDump("collection", "operation", models.MafraRetailPrice{}.TableName(), false)
}
