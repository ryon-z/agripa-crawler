package refiner

import (
	"errors"
	"fmt"
	"media_crawling/alarm"
	"media_crawling/config"
	"media_crawling/crawler"
	"media_crawling/models"
	"media_crawling/util"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// getBaseDate : 수출입 파일경로에서 기준일 추출
func getBaseDate(filePath string) string {
	_, file := filepath.Split(filePath)
	r := regexp.MustCompile(`^[A-Za-z]*_([0-9]{6}).csv$`)
	captured := r.FindStringSubmatch(file)

	if len(captured) != 2 {
		errorMessage := fmt.Sprintf("수출입 파일명이 잘못되었습니다. 경로: %s", filePath)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}

	var dateInts []int
	for _, dateStr := range []string{captured[1][:4], captured[1][4:6]} {
		dateInt, err := strconv.Atoi(dateStr)
		if err != nil {
			errorMessage := fmt.Sprintf("수출입 파일명의 날짜가 잘못되었습니다. 캡처 숫자: %s", captured)
			alarm.PostMessage("default", errorMessage)
			panic(err)
		}

		dateInts = append(dateInts, dateInt)
	}

	year := dateInts[0]
	month := dateInts[1]

	return fmt.Sprintf("%04d-%02d", year, month)
}

// RefineExportations : 수출 데이터 정제
func RefineExportations(filePath string) {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	var exportations []models.Exportation
	baseDate := getBaseDate(filePath)

	// 파일 읽음
	csvReader := util.ReadCsv(filePath)

	for _, row := range csvReader {
		weight := strings.ReplaceAll(row[4], ",", "")
		amount := strings.ReplaceAll(row[5], ",", "")
		exportations = append(exportations, models.Exportation{
			HskPrdlstCode: row[1], Weight: weight, Amount: amount, BaseDate: baseDate})
	}

	// exportations를 ImportDataToDB의 data 파라미터로 사용하기 위해
	// data 변수를 선언하고, exportations의 각 행을 할당한다.
	data := make([]interface{}, len(exportations))
	for index, row := range exportations {
		data[index] = row
	}
	// DB에 업로드
	crawler.ImportDataToDB(db, "exportation", "", data)

	fmt.Println(fmt.Sprintf("%s 수출 데이터 정제 및 업로드 완료", filePath))
}

// RefineImportations : 수입 데이터 정제
func RefineImportations(filePath string) {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	var importations []models.Importation
	baseDate := getBaseDate(filePath)

	// 파일 읽음
	csvReader := util.ReadCsv(filePath)

	for _, row := range csvReader {
		weight := strings.ReplaceAll(row[4], ",", "")
		amount := strings.ReplaceAll(row[5], ",", "")
		importations = append(importations, models.Importation{
			HskPrdlstCode: row[1], Weight: weight, Amount: amount, BaseDate: baseDate})
	}

	// importations를 ImportDataToDB의 data 파라미터로 사용하기 위해
	// data 변수를 선언하고, importations의 각 행을 할당한다.
	data := make([]interface{}, len(importations))
	for index, row := range importations {
		data[index] = row
	}
	// DB에 업로드
	crawler.ImportDataToDB(db, "importation", "", data)

	fmt.Println(fmt.Sprintf("%s 수입 데이터 정제 및 업로드 완료", filePath))
}

// checkTradeTypeCorrect : tradeType 타입 체크
func checkTradeTypeCorrect(tradeType string) {
	util.CheckCondition(
		!util.InArray(tradeType, []string{"exportation", "importation"}),
		"checkTradeTypeCorrect",
		fmt.Sprintf("허용되지 않은 tradeTpye: %s\n", tradeType))
}

// getAllTradeFilesPaths : 모든 수출입 파일 경로 획득
func getAllTradeFilesPaths(tradeType string) []string {
	checkTradeTypeCorrect(tradeType)
	dirPath := util.MakePath([]string{config.Conf.DataDirPath, tradeType})
	filesPaths, err := filepath.Glob(dirPath + "/*")
	util.CheckError(err, "getAllTradeFilesPaths")

	return filesPaths
}

// getTradeFileYearMonth : 입력 받은 수출입 파일명에서 날짜만 추출하여 int로 형변환 후 리턴
func getTradeFileYearMonth(tradeType string, fileName string) int {
	var result int
	checkTradeTypeCorrect(tradeType)
	fileName = strings.ReplaceAll(fileName, tradeType+"_", "")
	fileName = strings.ReplaceAll(fileName, ".csv", "")
	result, err := strconv.Atoi(fileName)
	util.CheckError(err, "getTradeFileYearMonth")

	return result
}

// getTradeFilesPaths : 특정 기간의 수출입 파일 경로 획득
func getTradeFilesPaths(tradeType string, start time.Time, end time.Time) []string {
	var result []string
	functionName := "getTradeFilesPaths"
	startYearMonth, err := strconv.Atoi(util.GetYearMonthString(start, ""))
	util.CheckError(err, functionName)
	endYearMonth, err := strconv.Atoi(util.GetYearMonthString(end, ""))
	util.CheckError(err, functionName)

	filePaths := getAllTradeFilesPaths(tradeType)
	for _, filePath := range filePaths {
		fileName := filepath.Base(filePath)
		fileYearMonth := getTradeFileYearMonth(tradeType, fileName)
		if fileYearMonth >= startYearMonth && fileYearMonth <= endYearMonth {
			result = append(result, filePath)
		}
	}

	return result
}

// RefineSpecificExportations : 특정 기간의 수출 데이터 정제 후 DB 업로드
func RefineSpecificExportations(start time.Time, end time.Time) {
	filesPaths := getTradeFilesPaths("exportation", start, end)
	for _, filePath := range filesPaths {
		RefineExportations(filePath)
	}
}

// RefineSpecificImportations : 특정 기간의 수입 데이터 정제 후 DB 업로드
func RefineSpecificImportations(start time.Time, end time.Time) {
	filesPaths := getTradeFilesPaths("importation", start, end)
	for _, filePath := range filesPaths {
		RefineImportations(filePath)
	}
}

// UpdateTradeFromImportation : 수입으로 부터 TRADE 테이블 업데이트
func UpdateTradeFromImportation(start time.Time, end time.Time) {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	startYearMonth := fmt.Sprintf("%04d-%02d", start.Year(), start.Month())
	endYearMonth := fmt.Sprintf("%04d-%02d", end.Year(), end.Month())
	tableName := models.Trade{}.TableName()
	importationTableName := models.Importation{}.TableName()

	sqlQuery := fmt.Sprintf(`
		INSERT IGNORE INTO %s (HskPrdlstCode, BaseDate, TradeType, Weight, Amount)
		SELECT A.HskPrdlstCode, BaseDate, "수입", Weight, Amount 
		FROM 
		(SELECT * FROM %s
		WHERE BaseDate >= "%s"
		AND BaseDate <= "%s") AS A
		JOIN (SELECT DISTINCT HskPrdlstCode from MAP_STD_HSK_ITEM) AS B
		ON A.HskPrdlstCode = B.HskPrdlstCode
	;`, tableName, importationTableName, startYearMonth, endYearMonth)
	crawler.RunQuery(db, sqlQuery)
}

// UpdateTradeFromExportation : 수출으로 부터 TRADE 테이블 업데이트
func UpdateTradeFromExportation(start time.Time, end time.Time) {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	startYearMonth := fmt.Sprintf("%04d-%02d", start.Year(), start.Month())
	endYearMonth := fmt.Sprintf("%04d-%02d", end.Year(), end.Month())
	tableName := models.Trade{}.TableName()
	exportationTableName := models.Exportation{}.TableName()

	sqlQuery := fmt.Sprintf(`
		INSERT IGNORE INTO %s (HskPrdlstCode, BaseDate, TradeType, Weight, Amount)
		SELECT A.HskPrdlstCode, BaseDate, "수출", Weight, Amount
		FROM
		(SELECT * FROM %s
		WHERE BaseDate >= "%s"
		AND BaseDate <= "%s") AS A
		JOIN (SELECT DISTINCT HskPrdlstCode FROM MAP_STD_HSK_ITEM) AS B
		ON A.HskPrdlstCode = B.HskPrdlstCode
	;`, tableName, exportationTableName, startYearMonth, endYearMonth)
	crawler.RunQuery(db, sqlQuery)
}
