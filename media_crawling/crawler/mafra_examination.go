package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"media_crawling/alarm"
	"media_crawling/config"
	"media_crawling/models"
	"media_crawling/util"
	"time"
)

func requestMafraExamin(startIndex int, endIndex int, date string) MafraExaminResponse {
	secret := config.Secret["mafra:secret"]
	metas := fmt.Sprintf("MAFRA_EXAMIN :: startIndex: %d, endIndex: %d", startIndex, endIndex)
	requestURL := "http://211.237.50.150:7080/openapi"
	completedURL := fmt.Sprintf("%s/%s/json/Grid_20160722000000000352_1/%d/%d?EXAMIN_DE=%s",
		requestURL, secret, startIndex, endIndex, date)

	var headers map[string]string
	headers = make(map[string]string)
	// 요청
	result := Request(completedURL, headers)
	fmt.Println("요청 성공")

	// json parsing
	response := MafraExaminResponse{}
	json.Unmarshal([]byte(result), &response)

	// mafra 요청 오류 체크
	resultCode := response.Grid201607220000000003521.Result.Code
	if resultCode != "INFO-000" {
		errorMessage := fmt.Sprintf("%s, 요청 실패, result: %s", metas, result)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}

	return response
}

func CrawlMafraExamination(start time.Time, end time.Time) {
	numOfRows := 999
	modelName := "mafraExamination"
	model := models.GetModel(modelName)
	rawTableName := model.TableName()

	db := GetDB("collection")
	defer CloseDB(db)

	for rd := util.GetRangeDate(start, end); ; {
		date := rd()
		if date.IsZero() {
			break
		}

		// 기준일 획득
		dateString := util.GetDateStringNoSep(date)

		// 테이블 명 재정의
		tableName := fmt.Sprintf("%s_%d", rawTableName, date.Year())
		CreateTableIfNotExists(db, tableName, rawTableName)

		// totalCounts 획득
		response := requestMafraExamin(1, 2, dateString)
		totalCounts := response.Grid201607220000000003521.TotalCnt

		fmt.Println("date", date)
		fmt.Println("totalCounts", totalCounts)

		for i := 1; i < totalCounts; i += numOfRows {
			response := requestMafraExamin(i, i+numOfRows-1, dateString)
			rows := response.Grid201607220000000003521.Row

			data := make([]interface{}, len(rows))
			for index, row := range rows {
				data[index] = row
			}
			// 여기까지 함
			fmt.Println("requesting is done")

			// DB에 저장
			ImportDataToDB(db, modelName, tableName, data)
			fmt.Println("importing data to db is done")

			time.Sleep(300 * time.Millisecond)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
