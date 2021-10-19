package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"media_crawling/alarm"
	"media_crawling/config"
	"media_crawling/util"
	"time"
)

// requestWeather : 일별 날씨 요청
func requestWeather(numOfRows int, pageNo int, startDate string, endDate string) WeatherResponse {
	secret := config.Secret["data.go.kr:secret"]
	metas := fmt.Sprintf("numOfRows: %d, pageNo: %d, startDate: %s, endDate: %s",
		numOfRows, pageNo, startDate, endDate)
	var requestURL string

	requestURL = "http://apis.data.go.kr/1360000/FmlandWthrInfoService/getDayStatistics"
	completedURL := fmt.Sprintf(`%s?serviceKey=%s&numOfRows=%d&pageNo=%d&ST_YMD=%s&ED_YMD=%s&AREA_ID=999999999&PA_CROP_SPE_ID=PA999999&dataType=json`,
		requestURL, secret, numOfRows, pageNo, startDate, endDate)

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
	weatherResponse := WeatherResponse{}
	json.Unmarshal([]byte(result), &weatherResponse)

	return weatherResponse
}

// getSubWeatherDates : 입력 받은 시작일과 종료일 사이의 년별로 하위 시작일과 종료일 리턴
func getSubWeatherDates(start time.Time, end time.Time) [][]string {
	var result [][]string
	startYear := start.Year()
	endYear := end.Year()
	startDate := util.GetDateStringNoSep(start)
	endDate := util.GetDateStringNoSep(end)

	for year := startYear; year <= endYear; year++ {
		if startYear == endYear {
			row := []string{startDate, endDate}
			result = append(result, row)
			break
		}

		thisYearEnd := fmt.Sprintf("%d1231", year)
		thisYearStart := fmt.Sprintf("%d0101", year)
		var row []string
		if year == startYear {
			row = []string{startDate, thisYearEnd}
		} else if year == endYear {
			row = []string{thisYearStart, endDate}
		} else {
			row = []string{thisYearStart, thisYearEnd}
		}
		result = append(result, row)
	}

	return result

}

// CrawlAllWeather : 모든 일별 날씨 수집
func CrawlAllWeather(start time.Time, end time.Time, logFilePath string) {
	fpLog, logger := util.GetFileLogger(logFilePath)
	defer util.CloseFileLogger(fpLog)

	db := GetDB("collection")
	defer CloseDB(db)

	numOfRows := 5000
	numRetrying := 5

	logger.Println("일별 날찌 수집 정상시작")

	for _, dates := range getSubWeatherDates(start, end) {
		subStartDate := dates[0]
		subEndDate := dates[1]
		fmt.Printf("subStartDate: %s, subEndDate: %s\n", subStartDate, subEndDate)

		weathers := requestWeather(1, 1, subStartDate, subEndDate)
		totalCount := weathers.Response.Body.TotalCount
		fmt.Println(totalCount)

		// 최대 페이지 수 획득
		maxPageNo := util.GetMaxPageNo(totalCount, numOfRows)
		fmt.Println("maxPageNo: ", maxPageNo)
		for pageNo := 1; pageNo < (maxPageNo + 1); pageNo++ {
			logger.Printf(
				"subStartDate: %s, subEndDate: %s, pageNo: %d 수집 시작\n",
				subStartDate, subEndDate, pageNo)

			weathers = requestWeather(numOfRows, pageNo, subStartDate, subEndDate)
			innerTotalCount := weathers.Response.Body.TotalCount
			fmt.Println("innerTotalCount: ", innerTotalCount)
			for i := 0; i < numRetrying; i++ {
				if innerTotalCount != 0 {
					break
				}

				weathers = requestWeather(numOfRows, pageNo, subStartDate, subEndDate)
				fmt.Printf("innser request %d 번째 재시도\n", i+1)
			}

			if innerTotalCount == 0 {
				metas := fmt.Sprintf("WEATHER :: numOfRows: %d, pageNo: %d, subStartDate: %s, subEndDate: %s",
					numOfRows, pageNo, subStartDate, subEndDate)
				errorMessage := fmt.Sprintf("%s, 내부 요청 실패", metas)
				alarm.PostMessage("default", errorMessage)
				panic(errors.New(errorMessage))
			}

			items := weathers.Response.Body.Items.Item
			fmt.Println("pageNo", pageNo)
			fmt.Println("len(items): ", len(items))
			data := make([]interface{}, len(items))
			for index, row := range items {
				data[index] = row
			}
			fmt.Println("requesting is done")
			// fmt.Println(data)

			// DB에 저장
			ImportDataToDB(db, "weather", "", data)
			fmt.Println("importing data to db is done")

			logger.Printf(
				"subStartDate: %s, subEndDate: %s, pageNo: %d 수집 종료\n",
				subStartDate, subEndDate, pageNo)
		}
	}

	logger.Println("일별 경락가격 통계정보 수집 정상종료")
}
