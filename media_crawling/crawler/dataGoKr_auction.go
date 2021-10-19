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

// CheckAuctionTypeValid : auctionType이 유효한지 체크
func CheckAuctionTypeValid(auctionType string) {
	if !util.InArray(auctionType, []string{"realtime", "adjusted"}) {
		errorMessage := fmt.Sprintf("유효하지 않은 auctionType. 입력 받은 auctionType: %s", auctionType)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}
}

// requestAuction : 경락 도매시장 유통 통계 요청
func requestAuction(auctionType string, numOfRows int, pageNo int, baseDate string, stdPrdlstCode string) AuctionResponse {
	secret := config.Secret["data.go.kr:secret"]
	metas := fmt.Sprintf("auctionType: %s, numOfRows: %d, pageNo: %d, baseDate: %s",
		auctionType, numOfRows, pageNo, baseDate)
	var requestURL string

	// request에 필요한 요소 구축
	CheckAuctionTypeValid(auctionType)
	if auctionType == "realtime" {
		requestURL = "http://apis.data.go.kr/B552895/StatsInfoService/getvRealtimeMktStatsInfo"
	} else {
		requestURL = "http://apis.data.go.kr/B552895/StatsInfoService/getDailyAdjStatsInfo"
	}

	completedURL := fmt.Sprintf("%s?serviceKey=%s&numOfRows=%d&pageNo=%d&delng_de=%s&_returnType=json",
		requestURL, secret, numOfRows, pageNo, baseDate)

	if stdPrdlstCode != "" {
		completedURL = fmt.Sprintf("%s&std_prdlst_code=%s", completedURL, stdPrdlstCode)
	}

	var headers map[string]string
	headers = make(map[string]string)
	headers["Accept-Charset"] = "UTF-8;q=1, ISO-8859-1;q=0"

	// 요청
	result := Request(completedURL, headers)
	isRetrying := IsRetryingDataGoKrResponse(result, metas)
	numRetrying := 5
	var i int
	for i = 0; i < numRetrying; i++ {
		if !isRetrying {
			break
		}

		time.Sleep(time.Second * 360)
		result := Request(completedURL, headers)
		isRetrying = IsRetryingDataGoKrResponse(result, metas)
		fmt.Printf("%d 번째 재시도\n", i+1)
	}
	if i == numRetrying {
		errorMessage := fmt.Sprintf("%s, 요청 실패", metas)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}
	fmt.Println("요청 성공")

	// json parsing
	realtimeAuctions := AuctionResponse{}
	json.Unmarshal([]byte(result), &realtimeAuctions)

	return realtimeAuctions
}

// CrawlAuction : 조사가격 품목코드의 일별 경락 도매시장 유통통계 수집 시작
func CrawlAuction(auctionType string, start time.Time, end time.Time) {
	var modelName string
	db := GetDB("collection")
	defer CloseDB(db)

	CheckAuctionTypeValid(auctionType)
	if auctionType == "realtime" {
		modelName = "realtimeAuction"
	} else {
		modelName = "adjustedAuction"
	}

	dateFormat := "20060102"
	numOfRows := 5000
	numRetrying := 5

	for rd := util.GetRangeDate(start, end); ; {
		date := rd()
		if date.IsZero() {
			break
		}

		for _, mapping := range models.GetStdSpciesCodesMappedExam(db) {
			time.Sleep(time.Second * 1)
			fmt.Println("mapping.StdSpciesCode:", mapping.StdSpciesCode)

			baseDate := date.Format(dateFormat)

			// totalCount 획득
			auctions := requestAuction(auctionType, 1, 1, baseDate, mapping.StdSpciesCode)
			totalCount := auctions.TotalCount

			// 최대 페이지 수 획득
			maxPageNo := util.GetMaxPageNo(totalCount, numOfRows)
			fmt.Println("maxPageNo: ", maxPageNo)
			for pageNo := 1; pageNo < (maxPageNo + 1); pageNo++ {
				// 기준일의 조사가격 품목코드 경락 도매시장 유통통계 요청
				auctions = requestAuction(auctionType, numOfRows, pageNo, baseDate, mapping.StdSpciesCode)
				fmt.Println("inner auctions.TotalCount: ", auctions.TotalCount)
				for i := 0; i < numRetrying; i++ {
					if auctions.TotalCount != 0 {
						break
					}

					auctions = requestAuction(auctionType, numOfRows, pageNo, baseDate, mapping.StdSpciesCode)
					fmt.Printf("innser request %d 번째 재시도\n", i+1)
				}
				if auctions.TotalCount == 0 {
					metas := fmt.Sprintf("auctionType: %s, numOfRows: %d, pageNo: %d, baseDate: %s",
						auctionType, numOfRows, pageNo, baseDate)
					errorMessage := fmt.Sprintf("%s, 내부 요청 실패", metas)
					alarm.PostMessage("default", errorMessage)
					panic(errors.New(errorMessage))
				}

				fmt.Println("pageNo", pageNo)
				fmt.Println("len(auctions.List): ", len(auctions.List))
				data := make([]interface{}, len(auctions.List))
				for index, row := range auctions.List {
					data[index] = row
				}
				fmt.Println("requesting is done")

				// DB에 저장
				ImportDataToDB(db, modelName, "", data)
				fmt.Println("importing data to db is done")
			}
		}
	}
}

// CrawlAllAuctions : 모든 일별 경락 도매시장 유통통계 수집(년도별 테이블)
func CrawlAllAuctions(auctionType string, start time.Time, end time.Time, logFilePath string) {
	var modelName string

	fpLog, logger := util.GetFileLogger(logFilePath)
	defer util.CloseFileLogger(fpLog)

	db := GetDB("collection")
	defer CloseDB(db)

	CheckAuctionTypeValid(auctionType)
	if auctionType == "realtime" {
		modelName = "realtimeAuction"
	} else {
		modelName = "adjustedAuction"
	}
	model := models.GetModel(modelName)
	rawTableName := model.TableName()

	dateFormat := "20060102"
	numOfRows := 5000
	numRetrying := 5

	logger.Println("일별 경락가격 통계정보 수집 정상시작")

	for rd := util.GetRangeDate(start, end); ; {
		date := rd()
		if date.IsZero() {
			break
		}

		baseDate := date.Format(dateFormat)

		// totalCount 획득
		auctions := requestAuction(auctionType, 1, 1, baseDate, "")
		totalCount := auctions.TotalCount
		// 최대 페이지 수 획득
		maxPageNo := util.GetMaxPageNo(totalCount, numOfRows)
		fmt.Println("maxPageNo: ", maxPageNo)

		// 딜레이
		// time.Sleep(time.Second * 1)

		// 테이블 명 재정의
		tableName := fmt.Sprintf("%s_%d", rawTableName, date.Year())
		CreateTableIfNotExists(db, tableName, "AGRI_ADJ_AUCTION")

		for pageNo := 1; pageNo < (maxPageNo + 1); pageNo++ {
			logger.Printf("baseDate: %s, pageNo: %d 수집 시작\n", baseDate, pageNo)

			// 기준일의 모든 품목코드 경락 도매시장 유통통계 요청
			auctions = requestAuction(auctionType, numOfRows, pageNo, baseDate, "")
			fmt.Println("inner auctions.TotalCount: ", auctions.TotalCount)
			for i := 0; i < numRetrying; i++ {
				if auctions.TotalCount != 0 {
					break
				}

				auctions = requestAuction(auctionType, numOfRows, pageNo, baseDate, "")
				fmt.Printf("innser request %d 번째 재시도\n", i+1)
			}
			if auctions.TotalCount == 0 {
				metas := fmt.Sprintf("auctionType: %s, numOfRows: %d, pageNo: %d, baseDate: %s",
					auctionType, numOfRows, pageNo, baseDate)
				errorMessage := fmt.Sprintf("%s, 내부 요청 실패", metas)
				alarm.PostMessage("default", errorMessage)
				panic(errors.New(errorMessage))
			}

			fmt.Println("pageNo", pageNo)
			fmt.Println("len(auctions.List): ", len(auctions.List))
			data := make([]interface{}, len(auctions.List))
			for index, row := range auctions.List {
				data[index] = row
			}

			// DB에 저장
			ImportDataToDB(db, modelName, tableName, data)
			fmt.Println("importing data to db is done")

			logger.Printf("baseDate: %s, pageNo: %d 수집 종료\n", baseDate, pageNo)
		}
	}

	logger.Println("일별 경락가격 통계정보 수집 정상종료")
}
