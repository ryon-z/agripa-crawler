package refiner

import (
	"fmt"
	"media_crawling/crawler"
	"media_crawling/models"
	"media_crawling/util"
	"time"
)

// UpdateMafraRetailPrice : MAFRA 소매 조사가격 테이블 업데이트
func UpdateMafraRetailPrice(start time.Time, end time.Time) {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	tableName := models.MafraRetailPrice{}.TableName()

	wherePhrases := GetDateWherePhrase("ExaminDe", "", start, end)
	for i, wherePhrase := range wherePhrases {
		year := start.Year() + i
		mafraExaminTableName := fmt.Sprintf("%s_%d", models.MafraExamination{}.TableName(), year)
		sqlQuery := fmt.Sprintf(`
			INSERT IGNORE INTO %s
			SELECT date_format(str_to_date(ExaminDe, '%%Y%%m%%d'),'%%Y-%%m-%%d'), ExaminPrdlstNm, ExaminPrdlstCode, ExaminSpciesNm, ExaminSpciesCode, ExaminUnitNm, ExaminUnit, ExaminGradNm, ExaminGradCode, min(TodayPric), max(TodayPric)
			FROM %s
			%s
			AND ExaminSeCode = "7"
			AND ExaminAreaCode = "1102"
			GROUP BY ExaminDe, ExaminPrdlstCode, ExaminSpciesCode, ExaminGradCode
		;`, tableName, mafraExaminTableName, wherePhrase)
		crawler.RunQuery(db, sqlQuery)
	}
}

// UpdateMafraWholePrice : MAFRA 도매시장 조사가격 테이블 업데이트
func UpdateMafraWholePrice(start time.Time, end time.Time) {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	tableName := models.MafraWholePrice{}.TableName()

	wherePhrases := GetDateWherePhrase("ExaminDe", "", start, end)
	for i, wherePhrase := range wherePhrases {
		year := start.Year() + i
		mafraExaminTableName := fmt.Sprintf("%s_%d", models.MafraExamination{}.TableName(), year)
		sqlQuery := fmt.Sprintf(`
			INSERT IGNORE INTO %s
			SELECT date_format(str_to_date(ExaminDe, '%%Y%%m%%d'),'%%Y-%%m-%%d'), ExaminPrdlstNm, ExaminPrdlstCode, ExaminSpciesNm, ExaminSpciesCode, ExaminUnitNm, ExaminUnit, ExaminGradNm, ExaminGradCode, min(TodayPric)
			FROM %s
			%s
			AND ExaminSeCode = "6"
			AND ExaminAreaCode = "1102"
			GROUP BY ExaminDe, ExaminPrdlstCode, ExaminSpciesCode, ExaminGradCode
		;`, tableName, mafraExaminTableName, wherePhrase)
		crawler.RunQuery(db, sqlQuery)
	}
}

func GetDateWherePhrase(dateColumnName string, dateStringSep string, start time.Time, end time.Time) []string {
	var result []string
	startYear := start.Year()
	endYear := end.Year()
	startDate := util.GetDateString(start, dateStringSep)
	endDate := util.GetDateString(end, dateStringSep)

	util.CheckCondition(startYear > endYear, "GetDateWherePhrase", "시작년도가 종료년도보다 큽니다.")

	if startYear == endYear {
		wherePhrase := fmt.Sprintf(`
			WHERE %s >= %s
			AND %s <= %s
		`, dateColumnName, startDate, dateColumnName, endDate)
		result = append(result, wherePhrase)

		return result
	}

	// startYear < endYear 일 경우
	for year := startYear; year <= endYear; year++ {
		if year < endYear {
			result = append(result, "WHERE 1")
		} else {
			wherePhrase := fmt.Sprintf(`
				WHERE %s >= %d0101
				AND %s <= %s
			`, dateColumnName, year, dateColumnName, endDate)
			result = append(result, wherePhrase)
		}
	}

	return result
}
