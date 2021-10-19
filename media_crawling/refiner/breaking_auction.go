package refiner

import (
	"fmt"
	"media_crawling/alarm"
	"media_crawling/crawler"
	"media_crawling/models"
)

// RemoveOldBreakingAuction : 오래된 실시간 경락 속보 삭제
func RemoveOldBreakingAuction(dateString string) {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	tableName := models.BreakingAuction{}.TableName()
	sqlQuery := fmt.Sprintf("DELETE FROM %s WHERE Bidtime < %s", tableName, dateString)

	// 데이터 삭제
	result := db.Exec(sqlQuery)
	if result.Error != nil {
		fmt.Println("RemoveOldBreakingAuction 에러 발생")

		alarm.PostMessage("default", result.Error.Error())
		panic(result.Error)
	}
}
