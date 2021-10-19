package crawler

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"media_crawling/alarm"
	"media_crawling/config"
	"strconv"
)

func requestGarakCode(numOfRows int, pageNo int) GarakCodeResponse {
	secret := config.Secret["data.go.kr:secret"]
	metas := fmt.Sprintf("GARAK_CODE ::  numOfRows: %d, pageNo: %d", numOfRows, pageNo)
	requestURL := "http://openapi.epis.or.kr/openapi/service/CodeListService/getPrdlstCodeList"
	completedURL := fmt.Sprintf("%s?secret=%s&numOfRows=%d&pageNo=%d",
		requestURL, secret, numOfRows, pageNo)

	var headers map[string]string
	headers = make(map[string]string)
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	// 요청
	result := RequestDataGoKr(completedURL, headers, metas)
	fmt.Println("요청 성공")

	// json parsing
	garakCodeResponse := GarakCodeResponse{}
	xml.Unmarshal([]byte(result), &garakCodeResponse)

	return garakCodeResponse
}

func requestWholesaleMarketCode(numOfRows int, pageNo int) WholesaleMarketCodeResponse {
	secret := config.Secret["data.go.kr:secret"]
	metas := fmt.Sprintf("WHOLESALE_MARKET_CODE :: numOfRows: %d, pageNo: %d", numOfRows, pageNo)
	requestURL := "http://openapi.epis.or.kr/openapi/service/CodeListService/getWltCodeList"
	completedURL := fmt.Sprintf("%s?secret=%s&numOfRows=%d&pageNo=%d",
		requestURL, secret, numOfRows, pageNo)

	var headers map[string]string
	headers = make(map[string]string)
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	// 요청
	result := RequestDataGoKr(completedURL, headers, metas)
	fmt.Println("요청 성공")

	// json parsing
	wholesaleMarketCodeResponse := WholesaleMarketCodeResponse{}
	xml.Unmarshal([]byte(result), &wholesaleMarketCodeResponse)

	return wholesaleMarketCodeResponse
}

func requestWholesaleMarketCoCode(numOfRows int, pageNo int) WholesaleMarketCoCodeResponse {
	secret := config.Secret["data.go.kr:secret"]
	metas := fmt.Sprintf("WHOLESALE_MARKET_CO_CODE :: numOfRows: %d, pageNo: %d", numOfRows, pageNo)
	requestURL := "http://openapi.epis.or.kr/openapi/service/CodeListService/getWltprCodeList"
	completedURL := fmt.Sprintf("%s?secret=%s&numOfRows=%d&pageNo=%d",
		requestURL, secret, numOfRows, pageNo)

	var headers map[string]string
	headers = make(map[string]string)
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	// 요청
	result := RequestDataGoKr(completedURL, headers, metas)
	fmt.Println("요청 성공")

	// json parsing
	wholesaleMarketCoCodeResponse := WholesaleMarketCoCodeResponse{}
	xml.Unmarshal([]byte(result), &wholesaleMarketCoCodeResponse)

	return wholesaleMarketCoCodeResponse
}

func requestStdGradeCode(numOfRows int, pageNo int) StdGradeCodeResponse {
	secret := config.Secret["data.go.kr:secret"]
	metas := fmt.Sprintf("STD_GRADE_CODE :: numOfRows: %d, pageNo: %d", numOfRows, pageNo)
	requestURL := "http://openapi.epis.or.kr/openapi/service/CodeListService/getGradCodeList"
	completedURL := fmt.Sprintf("%s?secret=%s&numOfRows=%d&pageNo=%d",
		requestURL, secret, numOfRows, pageNo)

	var headers map[string]string
	headers = make(map[string]string)
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	// 요청
	result := RequestDataGoKr(completedURL, headers, metas)
	fmt.Println("요청 성공")

	// json parsing
	stdgradeCodeResponse := StdGradeCodeResponse{}
	xml.Unmarshal([]byte(result), &stdgradeCodeResponse)

	return stdgradeCodeResponse
}

func requestStdUnitCode(numOfRows int, pageNo int) StdUnitCodeResponse {
	secret := config.Secret["data.go.kr:secret"]
	metas := fmt.Sprintf("STD_UNIT_CODE :: numOfRows: %d, pageNo: %d", numOfRows, pageNo)
	requestURL := "http://openapi.epis.or.kr/openapi/service/CodeListService/getUnitCodeList"
	completedURL := fmt.Sprintf("%s?secret=%s&numOfRows=%d&pageNo=%d",
		requestURL, secret, numOfRows, pageNo)

	var headers map[string]string
	headers = make(map[string]string)
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	// 요청
	result := RequestDataGoKr(completedURL, headers, metas)
	fmt.Println("요청 성공")

	// json parsing
	stdUnitCodeResponse := StdUnitCodeResponse{}
	xml.Unmarshal([]byte(result), &stdUnitCodeResponse)

	return stdUnitCodeResponse
}

func requestPlaceOriginCode(numOfRows int, pageNo int) PlaceOriginCodeResponse {
	secret := config.Secret["data.go.kr:secret"]
	metas := fmt.Sprintf("PLACE_ORIGIN_CODE :: numOfRows: %d, pageNo: %d", numOfRows, pageNo)
	requestURL := "http://openapi.epis.or.kr/openapi/service/CodeListService/getMtcCodeList"
	completedURL := fmt.Sprintf("%s?secret=%s&numOfRows=%d&pageNo=%d",
		requestURL, secret, numOfRows, pageNo)

	var headers map[string]string
	headers = make(map[string]string)
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	// 요청
	result := RequestDataGoKr(completedURL, headers, metas)
	fmt.Println("요청 성공")

	// json parsing
	placeOriginCodeResponse := PlaceOriginCodeResponse{}
	xml.Unmarshal([]byte(result), &placeOriginCodeResponse)

	return placeOriginCodeResponse
}

// requestStdSpeciesCode : 표준품종코드 요청
func requestStdSpeciesCode(startIndex int, endIndex int) StdSpeciesResponse {
	secret := config.Secret["mafra:secret"]
	metas := fmt.Sprintf("STD_SPECIES_CODE :: startIndex: %d, endIndex: %d", startIndex, endIndex)
	requestURL := "http://211.237.50.150:7080/openapi"
	completedURL := fmt.Sprintf("%s/%s/json/Grid_20141221000000000120_1/%d/%d",
		requestURL, secret, startIndex, endIndex)

	var headers map[string]string
	headers = make(map[string]string)
	// 요청
	result := Request(completedURL, headers)
	fmt.Println("요청 성공")

	// json parsing
	stdSpeciesResponse := StdSpeciesResponse{}
	json.Unmarshal([]byte(result), &stdSpeciesResponse)

	// mafra 요청 오류 체크
	resultCode := stdSpeciesResponse.Grid201412210000000001201.Result.Code
	if resultCode != "INFO-000" {
		errorMessage := fmt.Sprintf("%s, 요청 실패, result: %s", metas, result)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}

	return stdSpeciesResponse
}

func CrawlAllGarakCodes() {
	db := GetDB("collection")
	defer CloseDB(db)

	response := requestGarakCode(1, 1)
	totalCount, err := strconv.Atoi(response.Body.TotalCount)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	response = requestGarakCode(totalCount, 1)
	items := response.Body.Items.Item
	data := make([]interface{}, len(items))
	for index, row := range items {
		data[index] = row
	}

	// DB에 저장
	ImportDataToDB(db, "garakCode", "", data)
	fmt.Println("importing data to db is done")
}

func CrawlAllWholesaleMarketCodes() {
	db := GetDB("collection")
	defer CloseDB(db)

	response := requestWholesaleMarketCode(1, 1)
	totalCount, err := strconv.Atoi(response.Body.TotalCount)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	response = requestWholesaleMarketCode(totalCount, 1)
	items := response.Body.Items.Item
	data := make([]interface{}, len(items))
	for index, row := range items {
		data[index] = row
	}

	// DB에 저장
	ImportDataToDB(db, "wholesaleMarketCode", "", data)
	fmt.Println("importing data to db is done")
}

func CrawlAllWholesaleMarketCoCodes() {
	db := GetDB("collection")
	defer CloseDB(db)

	response := requestWholesaleMarketCoCode(1, 1)
	totalCount, err := strconv.Atoi(response.Body.TotalCount)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	response = requestWholesaleMarketCoCode(totalCount, 1)

	items := response.Body.Items.Item
	data := make([]interface{}, len(items))
	for index, row := range items {
		data[index] = row
	}

	// DB에 저장
	ImportDataToDB(db, "wholesaleMarketCoCode", "", data)
	fmt.Println("importing data to db is done")
}

func CrawlAllStdGradeCode() {
	db := GetDB("collection")
	defer CloseDB(db)

	response := requestStdGradeCode(1, 1)
	totalCount, err := strconv.Atoi(response.Body.TotalCount)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	response = requestStdGradeCode(totalCount, 1)

	items := response.Body.Items.Item
	data := make([]interface{}, len(items))
	for index, row := range items {
		data[index] = row
	}

	// DB에 저장
	ImportDataToDB(db, "stdGradeCode", "", data)
	fmt.Println("importing data to db is done")
}

func CrawlAllStdUnitCode() {
	db := GetDB("collection")
	defer CloseDB(db)

	response := requestStdUnitCode(1, 1)
	totalCount, err := strconv.Atoi(response.Body.TotalCount)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	response = requestStdUnitCode(totalCount, 1)

	items := response.Body.Items.Item
	data := make([]interface{}, len(items))
	for index, row := range items {
		data[index] = row
	}

	// DB에 저장
	ImportDataToDB(db, "stdUnitCode", "", data)
	fmt.Println("importing data to db is done")
}

func CrawlAllPlaceOriginCode() {
	db := GetDB("collection")
	defer CloseDB(db)

	response := requestPlaceOriginCode(1, 1)
	totalCount, err := strconv.Atoi(response.Body.TotalCount)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	response = requestPlaceOriginCode(totalCount, 1)
	items := response.Body.Items.Item
	data := make([]interface{}, len(items))
	for index, row := range items {
		data[index] = row
	}

	// DB에 저장
	ImportDataToDB(db, "placeOriginCode", "", data)
	fmt.Println("importing data to db is done")
}

func CrawlAllStdSpeciesCode() {
	db := GetDB("collection")
	defer CloseDB(db)

	maxDiff := 999
	response := requestStdSpeciesCode(1, 2)
	totalCounts := response.Grid201412210000000001201.TotalCnt

	for i := 1; i <= totalCounts; i += maxDiff {
		response = requestStdSpeciesCode(i, i+maxDiff)
		items := response.Grid201412210000000001201.Row
		data := make([]interface{}, len(items))
		for index, row := range items {
			data[index] = row
		}

		// DB에 저장
		ImportDataToDB(db, "stdSpeciesCode", "", data)
		fmt.Println("importing data to db is done")

	}
}
