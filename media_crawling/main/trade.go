package main

import (
	"fmt"
	"media_crawling/config"
	"media_crawling/crawler"
	"media_crawling/models"
	"media_crawling/refiner"
	"media_crawling/util"
	"time"
)

func main() {
	conf := config.Conf
	loc := conf.Timezone

	// 임의 기간의 수출입 데이터 정제
	// (데이터 업로드가 수동이라 자동화 불가)
	start := time.Date(2010, 1, 1, 0, 0, 0, 0, loc)
	end := time.Date(2020, 8, 1, 0, 0, 0, 0, loc)

	startYearMonth := util.GetDateString(start, "-")[:7]
	endYearMonth := util.GetDateString(end, "-")[:7]
	dateColumnName := "BaseDate"

	refiner.RefineSpecificExportations(start, end)
	refiner.RefineSpecificImportations(start, end)
	refiner.UpdateTradeFromImportation(start, end)
	refiner.UpdateTradeFromExportation(start, end)
	sqlQuery := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE %s >= "%s"
		AND %s <= "%s"
		;`, models.Trade{}.TableName(),
		dateColumnName, startYearMonth, dateColumnName, endYearMonth)
	crawler.MoveDataUsingGorm("collection", "operation", "trade", sqlQuery)

	// 덤프 명령어
	// crawler.MoveDataUsingDump("collection", "operation", models.Trade{}.TableName(), false)
}
