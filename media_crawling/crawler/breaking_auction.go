package crawler

import (
	"encoding/xml"
	"errors"
	"fmt"
	"media_crawling/alarm"
	"media_crawling/config"
	"media_crawling/models"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
)

func requestBreakingAuction(numOfRows int, pageNo int, date string) BreakingAuctionResponse {
	secret := config.Secret["data.go.kr:secret"]
	metas := fmt.Sprintf("BREAKING_AUCTION :: numOfRows: %d, pageNo: %d", numOfRows, pageNo)
	requestURL := "http://openapi.epis.or.kr/openapi/service/RltmAucBrknewsService/getWltRltmAucBrknewsList"
	completedURL := fmt.Sprintf("%s?secret=%s&numOfRows=%d&pageNo=%d&dates=%s",
		requestURL, secret, numOfRows, pageNo, date)

	var headers map[string]string
	headers = make(map[string]string)
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	// 요청
	result := RequestDataGoKr(completedURL, headers, metas)
	fmt.Println("요청 성공")

	// json parsing
	breakingAuctionResponse := BreakingAuctionResponse{}
	xml.Unmarshal([]byte(result), &breakingAuctionResponse)

	return breakingAuctionResponse
}

func refineBidTime(date string, bidTime string) string {
	time := strings.Split(bidTime, " ")

	return fmt.Sprintf("%s-%s-%s %s:%s",
		date[:4], date[4:6], date[6:], time[1], time[2])
}

// getWholesaleMarketCoCodeMap : 도매시장 법인 코드 맵 획득
func getWholesaleMarketCoCodeMap(db *gorm.DB) map[string]string {
	var result map[string]string
	result = make(map[string]string)
	var wholesaleMarketCoCodes []models.WholesaleMarketCoCode
	tableName := models.WholesaleMarketCoCode{}.TableName()
	sqlQuery := fmt.Sprintf("SELECT * FROM %s;", tableName)
	db.Raw(sqlQuery).Scan(&wholesaleMarketCoCodes)

	for _, wholesaleMarketCoCode := range wholesaleMarketCoCodes {
		key := wholesaleMarketCoCode.Coname
		value := wholesaleMarketCoCode.Cocode
		result[key] = value
	}

	return result
}

// getStdGradeCodeMap : 표준 등급 코드 맵 획득
func getStdGradeCodeMap(db *gorm.DB) map[string]string {
	var result map[string]string
	result = make(map[string]string)
	var stdGradeCodes []models.StdGradeCode
	tableName := models.StdGradeCode{}.TableName()
	sqlQuery := fmt.Sprintf("SELECT * FROM %s;", tableName)
	db.Raw(sqlQuery).Scan(&stdGradeCodes)

	for _, stdGradeCode := range stdGradeCodes {
		key := stdGradeCode.Gradename
		value := stdGradeCode.Gradecode
		result[key] = value
	}

	return result
}

// getWholesaleMarketCodeMap : 도매시장 코드 맵 획득
func getWholesaleMarketCodeMap(db *gorm.DB) map[string]string {
	var result map[string]string
	result = make(map[string]string)
	var wholesaleMarketCodes []models.WholesaleMarketCode
	tableName := models.WholesaleMarketCode{}.TableName()
	sqlQuery := fmt.Sprintf("SELECT * FROM %s;", tableName)
	db.Raw(sqlQuery).Scan(&wholesaleMarketCodes)

	for _, wholesaleMarketCode := range wholesaleMarketCodes {
		key := wholesaleMarketCode.Marketnm
		value := wholesaleMarketCode.Marketco
		result[key] = value
	}

	return result
}

// getStdUnitCodeMap : 표준 단위 코드 맵 획득
func getStdUnitCodeMap(db *gorm.DB) map[string]string {
	var result map[string]string
	result = make(map[string]string)
	var stdUnitCodes []models.StdUnitCode
	tableName := models.StdUnitCode{}.TableName()
	sqlQuery := fmt.Sprintf("SELECT * FROM %s;", tableName)
	db.Raw(sqlQuery).Scan(&stdUnitCodes)

	space := regexp.MustCompile(`\s+`)
	for _, stdUnitCode := range stdUnitCodes {
		key := space.ReplaceAllString(stdUnitCode.Unitname, " ")
		value := stdUnitCode.Unitcode
		result[key] = value
	}

	return result
}

// getStdSpeciesCodeMap : 표준 품종 코드 맵 획득
func getStdSpeciesCodeMap(db *gorm.DB) (map[string]string, map[string]string) {
	var specieseResult, itemResult map[string]string
	specieseResult = make(map[string]string)
	itemResult = make(map[string]string)
	var stdSpeciesCodes []models.StdSpeciesCode
	tableName := models.StdSpeciesCode{}.TableName()
	sqlQuery := fmt.Sprintf("SELECT * FROM %s;", tableName)
	db.Raw(sqlQuery).Scan(&stdSpeciesCodes)

	for _, stdSpeciesCode := range stdSpeciesCodes {
		key := stdSpeciesCode.MClassName
		value := stdSpeciesCode.MClassCode
		specieseResult[key] = value

		key = stdSpeciesCode.SClassName
		value = stdSpeciesCode.SClassCode
		itemResult[key] = value
	}

	return itemResult, specieseResult
}

// getPlaceOriginCodes : 산지 코드 구조체 획득
func getPlaceOriginCodes(db *gorm.DB) []models.PlaceOriginCode {
	var placeOriginCodes []models.PlaceOriginCode
	tableName := models.PlaceOriginCode{}.TableName()
	sqlQuery := fmt.Sprintf("SELECT * FROM %s;", tableName)
	db.Raw(sqlQuery).Scan(&placeOriginCodes)

	return placeOriginCodes
}

func getSomeCode(name string, codeMap map[string]string) string {
	result := "^"
	if correctCode, ok := codeMap[name]; ok {
		result = correctCode
	}

	return result
}

func CrawlBreakingAuction(date string) {
	// date 에러 처리
	if len(date) != 8 {
		errorMessage := fmt.Sprint("BREAKING_AUCTION :: 입력 날짜 길이가 8이 아닙니다.")
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}

	db := GetDB("collection")
	defer CloseDB(db)

	numOfRows := 500
	response := requestBreakingAuction(numOfRows, 1, date)
	items := response.Body.Items.Item
	data := make([]interface{}, len(items))

	wholesaleMarketcoCodeMap := getWholesaleMarketCoCodeMap(db)
	stdGradeCodeMap := getStdGradeCodeMap(db)
	wholesaleMarketCodeMap := getWholesaleMarketCodeMap(db)
	stdUnitCodeMap := getStdUnitCodeMap(db)
	speciesCodeMap, itemCodeMap := getStdSpeciesCodeMap(db)
	placeOriginCodes := getPlaceOriginCodes(db)

	for index, row := range items {
		refinedRow := models.BreakingAuction{}
		refinedRow.Bidtime = refineBidTime(date, row.Bidtime)
		refinedRow.Sanji = strings.ReplaceAll(row.Sanji, "-", " ")
		refinedRow.Unitname = strings.TrimSpace(row.Unitname)
		refinedRow.Coname = row.Coname
		refinedRow.Gradename = row.Gradename
		refinedRow.Marketname = row.Marketname
		refinedRow.Mclassname = row.Mclassname
		refinedRow.Price = row.Price
		refinedRow.Sclassname = row.Sclassname
		refinedRow.Tradeamt = row.Tradeamt
		refinedRow.Chulagtnm = row.Chulagtnm
		refinedRow.Cocode = getSomeCode(row.Coname, wholesaleMarketcoCodeMap)
		refinedRow.Gradecode = getSomeCode(row.Gradename, stdGradeCodeMap)
		refinedRow.Marketco = getSomeCode(row.Marketname, wholesaleMarketCodeMap)
		refinedRow.MclassCode = getSomeCode(row.Mclassname, itemCodeMap)

		refinedSclassName := strings.Replace(row.Sclassname, row.Mclassname+" ", "", 1)
		refinedRow.SClassCode = getSomeCode(refinedSclassName, speciesCodeMap)
		if refinedRow.SClassCode == "^" {
			refinedRow.SClassCode = getSomeCode(row.Sclassname, speciesCodeMap)
		}

		zipCode := "^"
	out:
		for _, placeOriginCode := range placeOriginCodes {
			candidateSanjis := []string{
				placeOriginCode.Sido + " " + placeOriginCode.Sigun,
				placeOriginCode.Sido + " " + placeOriginCode.Sigun + " " + placeOriginCode.Dong,
				placeOriginCode.Sigun,
				placeOriginCode.Sido,
				placeOriginCode.Dong,
			}

			for _, candidateSangi := range candidateSanjis {
				if refinedRow.Sanji == candidateSangi {
					zipCode = placeOriginCode.Zipcode
					break out
				}
			}
		}
		refinedRow.Zipcode = zipCode

		refinedRow.Unitamt = "0"
		refinedRow.Unitcode = "^"
		space := regexp.MustCompile(`\s+`)
		refinedUnitName := space.ReplaceAllString(row.Unitname, " ")
		re, _ := regexp.Compile(`^([0-9]*\.?[0-9]*)(.*)`)
		subMatches := re.FindStringSubmatch(refinedUnitName)
		candidateUnitAmt := subMatches[1]
		candidateUnitName := subMatches[2]
		if candidateUnitAmt != "" {
			refinedRow.Unitamt = subMatches[1]
		}
		if candidateUnitName != "" {
			refinedRow.Unitcode = getSomeCode(candidateUnitName, stdUnitCodeMap)
		}

		data[index] = refinedRow
	}

	// DB에 저장
	ImportDataToDB(db, "breakingAuction", "", data)
	fmt.Println("importing data to db is done")
}

func CrawlRawBreakingAuction(date string) {
	// date 에러 처리
	if len(date) != 8 {
		errorMessage := fmt.Sprint("BREAKING_AUCTION :: 입력 날짜 길이가 8이 아닙니다.")
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}

	db := GetDB("collection")
	defer CloseDB(db)

	numOfRows := 300
	response := requestBreakingAuction(numOfRows, 1, date)
	items := response.Body.Items.Item
	data := make([]interface{}, len(items))
	for index, row := range items {
		row.Bidtime = refineBidTime(date, row.Bidtime)
		row.Sanji = strings.ReplaceAll(row.Sanji, "-", " ")
		row.Unitname = strings.TrimSpace(row.Unitname)
		data[index] = row
	}

	// DB에 저장
	ImportDataToDB(db, "rawBreakingAuction", "", data)
	fmt.Println("importing data to db is done")
}
