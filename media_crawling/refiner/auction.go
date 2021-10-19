package refiner

import (
	"fmt"
	"media_crawling/crawler"
	"media_crawling/models"
	"media_crawling/util"
	"time"

	"github.com/jinzhu/gorm"
)

// GetAuctionPeriodWherePhrases : 경락가격 기간 조건 구 획득
func GetAuctionPeriodWherePhrases(start time.Time, end time.Time) []string {
	var result []string
	startYear := start.Year()
	endYear := end.Year()
	startPhrase := util.GetDateStringNoSep(start)
	endPhrase := util.GetDateStringNoSep(end)

	for year := startYear; year <= endYear; year++ {
		if startYear == endYear {
			row := fmt.Sprintf(" WHERE DelngDe >= \"%s\" AND DelngDe <= \"%s\" ", startPhrase, endPhrase)
			result = append(result, row)
			break
		}

		if year == startYear {
			row := fmt.Sprintf(" WHERE DelngDe >= \"%s\" AND DelngDe <= \"%d1231\" ", startPhrase, startYear)
			result = append(result, row)
		} else if year == endYear {
			row := fmt.Sprintf(" WHERE DelngDe >= \"%d0101\" AND DelngDe <= \"%s\" ", endYear, endPhrase)
			result = append(result, row)
		} else {
			result = append(result, " ")
		}
	}

	return result
}

// refineRealtimeAuction : DATA.GO.KR 일별 실시간 경락가격 통계정보 정제 및 DB import
func refineRealtimeAuction(db *gorm.DB, rawTableName string, tableName string, wherePhrase string) []models.RealtimeAuction {
	var auctions []models.RealtimeAuction
	sqlQuery := fmt.Sprintf(`
	INSERT IGNORE INTO %s
	SELECT A.* 
	FROM %s AS A
	INNER JOIN (SELECT DISTINCT StdSpciesCode from MAP_STD_EXAM_ITEM) AS B
	ON A.StdPrdlstCode = B.StdSpciesCode
	%s
	;
	`, rawTableName, tableName, wherePhrase)
	fmt.Printf("refineRealtimeAuction sqlQuery: %s\n", sqlQuery)
	db.Raw(sqlQuery).Scan(&auctions)

	return auctions
}

// refineAdjustedAuction : DATA.GO.KR 일별 정산 경락가격 통계정보 정제 및 DB import
func refineAdjustedAuction(db *gorm.DB, rawTableName string, tableName string, wherePhrase string) []models.AdjustedAuction {
	var auctions []models.AdjustedAuction
	sqlQuery := fmt.Sprintf(`
	INSERT IGNORE INTO %s
	SELECT A.* 
	FROM %s AS A
	INNER JOIN (SELECT DISTINCT StdSpciesCode from MAP_STD_EXAM_ITEM) AS B
	ON A.StdPrdlstCode = B.StdSpciesCode
	%s
	;
	`, rawTableName, tableName, wherePhrase)
	fmt.Printf("refineAdjustedAuction sqlQuery: %s\n", sqlQuery)
	db.Raw(sqlQuery).Scan(&auctions)

	return auctions
}

// RefineAcution : 경락가격 정제 후 DB import
func RefineAcution(auctionType string, start time.Time, end time.Time) {
	var modelName string
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	crawler.CheckAuctionTypeValid(auctionType)
	if auctionType == "realtime" {
		modelName = "realtimeAuction"
	} else {
		modelName = "adjustedAuction"
	}

	model := models.GetModel(modelName)
	rawTableName := model.TableName()

	wherePhrases := GetAuctionPeriodWherePhrases(start, end)
	years := util.GetYearsInRange(start, end)

	for i, phrase := range wherePhrases {
		tableName := fmt.Sprintf("%s_%d", rawTableName, years[i])
		if auctionType == "realtime" {
			refineRealtimeAuction(db, rawTableName, tableName, phrase)
		} else {
			refineAdjustedAuction(db, rawTableName, tableName, phrase)
		}
	}
}
