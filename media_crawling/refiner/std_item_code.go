package refiner

import (
	"fmt"
	"media_crawling/crawler"
	"media_crawling/models"
	"media_crawling/util"
	"regexp"
	"strings"
)

// SetStdItemKeword : 표준품목코드 키워드 세팅
func SetStdItemKeword() {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	tableName := models.StdItemKeyword{}.TableName()
	rawtableName := models.StdSpeciesCode{}.TableName()

	var stdSpeciesCodes []models.StdSpeciesCode
	var stdItemKeywords []models.StdItemKeyword

	// ItemCode 갱신
	sqlQuery := fmt.Sprintf(`
		DELETE FROM %s
		WHERE ItemCode NOT IN (SELECT DISTINCT MClassCode FROM %s)
	;`, tableName, rawtableName)
	crawler.RunQuery(db, sqlQuery)

	// 표준품목명 삽입
	sqlQuery = fmt.Sprintf(`
		INSERT IGNORE INTO %s (ItemCode, ExposedKeyWord, KeyWord, Priority)
		SELECT DISTINCT MClassCode, MClassName, MClassName, 1 FROM %s
	;`, tableName, rawtableName)
	crawler.RunQuery(db, sqlQuery)

	// StdSpeciesCode 조회
	sqlQuery = fmt.Sprintf(`
		SELECT * FROM %s
	;`, rawtableName)
	db.Raw(sqlQuery).Scan(&stdSpeciesCodes)

	// 품목명(품종명) 형태로 정제
	removeSpecialReg, err := regexp.Compile("[^a-zA-Z0-9ㄱ-ㅎㅏ-ㅣ가-힣]+")
	util.CheckError(err, "TestItemKeyword")
	for _, stdSpeciesCode := range stdSpeciesCodes {
		itemCode := stdSpeciesCode.MClassCode
		itemName := stdSpeciesCode.MClassName
		speciesName := stdSpeciesCode.SClassName
		if speciesName == "기타" {
			continue
		}

		removeFirstReg, err := regexp.Compile(fmt.Sprintf("^[%s]+", itemName))
		util.CheckError(err, "TestItemKeyword")

		replacedSpeciesName := removeFirstReg.ReplaceAllString(speciesName, "")
		replacedSpeciesName = removeSpecialReg.ReplaceAllString(replacedSpeciesName, "")

		var refined string
		if replacedSpeciesName == "" {
			refined = itemName
		} else {
			refined = fmt.Sprintf("%s(%s)", itemName, replacedSpeciesName)
		}

		var row models.StdItemKeyword
		row.ItemCode = itemCode
		row.ExposedKeyword = strings.TrimSpace(itemName)
		row.Keyword = strings.TrimSpace(refined)
		row.IsDisplay = "1"
		row.NumSearch = "0"
		row.Priority = "0"

		stdItemKeywords = append(stdItemKeywords, row)
	}

	// DB에 적재
	data := make([]interface{}, len(stdItemKeywords))
	for index, row := range stdItemKeywords {
		data[index] = row
	}
	crawler.ImportDataToDB(db, "stdItemKeyword", "", data)
}

// SetItemMapping : ItemMapping 테이블 세팅
func SetItemMapping() {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	tableName := models.ItemMapping{}.TableName()
	rawTableName := models.StdSpeciesCode{}.TableName()

	// ITEM_MAPPING 테이블 초기화
	sqlQuery := fmt.Sprintf(`DELETE FROM %s;`, tableName)
	crawler.RunQuery(db, sqlQuery)

	sqlQuery = fmt.Sprintf(`
		INSERT IGNORE INTO %s
		SELECT MClassCode, 
			CASE WHEN ExaminPrdlstCode is null THEN "non" else ExaminPrdlstCode END, 
			CASE WHEN HskPrdlstCode is null THEN "non" else HskPrdlstCode END 
		FROM
		(SELECT DISTINCT MClassCode FROM %s) AS A
		LEFT JOIN (select distinct StdPrdlstCode, ExaminPrdlstCode FROM MAP_STD_EXAM_ITEM) AS B
		ON A.MclassCode = B.StdPrdlstCode
		LEFT JOIN (select distinct StdPrdlstCode, HskPrdlstCode FROM MAP_STD_HSK_ITEM WHERE HskPrdlstCode != "") AS C
		ON A.MclassCode = C.StdPrdlstCode
	;`, tableName, rawTableName)
	crawler.RunQuery(db, sqlQuery)
}

// SetItemCode : 한 품목코드 당 하나의 품목명만 갖는 ItemCode 테이블 세팅
func SetItemCode() {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	tableName := models.StdItemCode{}.TableName()

	// STD_ITEM_CODE 테이블 초기화
	sqlQuery := fmt.Sprintf(`DELETE FROM %s;`, tableName)
	crawler.RunQuery(db, sqlQuery)

	sqlQuery = fmt.Sprintf(`
		INSERT IGNORE INTO %s
		SELECT MClassCode, MClassName FROM STD_SPECIES_CODE
		GROUP BY MClassCode
	;`, tableName)
	crawler.RunQuery(db, sqlQuery)

}
