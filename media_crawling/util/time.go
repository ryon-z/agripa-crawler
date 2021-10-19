package util

import (
	"fmt"
	"media_crawling/alarm"
	"media_crawling/config"
	"strings"
	"time"
)

// GetKoreanDateTime : datetimeFormat에 맞게 시간을 변환합니다. UTC시간이라면 한국시간으로 변환합니다.
// (timeType 허용값: "RFC1123Z", "RFC3339")
func GetKoreanDateTime(datetime string, timeType string, isUTC bool, datetimeFormat string) string {
	var t time.Time
	var koreanT time.Time
	allowedTimeType := []string{"RFC1123Z", "RFC3339"}

	if timeType == "RFC1123Z" {
		t, _ = time.Parse(time.RFC1123Z, datetime)
	} else if timeType == "RFC3339" {
		t, _ = time.Parse(time.RFC3339, datetime)
	} else {
		errorMessage := fmt.Sprintf(
			"잘못된 timeType 입니다. timetype: %s, allowedTimeType: %s",
			timeType, strings.Join(allowedTimeType, ", "))

		alarm.PostMessage("default", errorMessage)
		panic(errorMessage)
	}

	if isUTC {
		koreanT = t.Add(time.Hour * 9)
	} else {
		koreanT = t
	}

	koreanDatetime := koreanT.Format(datetimeFormat)

	return koreanDatetime
}

// GetTimeForMonthsAgo : 1 달 전 시간을 얻습니다(RFC1123Z유형의 문자열로 리턴).
func GetTimeForMonthsAgo(months int) string {
	t := time.Now().In(config.Conf.Timezone)

	return t.AddDate(0, months*-1, 0).Format(time.RFC3339)
}

// GetRangeDate : 입력 받은 기간 사이의 날짜를 모두 리턴.
// Reference by : https://stackoverflow.com/questions/50982524/how-to-gracefully-iterate-a-date-range-in-go
func GetRangeDate(start, end time.Time) func() time.Time {
	loc := config.Conf.Timezone
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, loc)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, loc)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		fmt.Println(date)
		return date
	}
}

// GetDateStringNoSep : 구분자가 없는 날짜 문자열 리턴
func GetDateStringNoSep(date time.Time) string {
	return fmt.Sprintf("%02d%02d%02d", date.Year(), date.Month(), date.Day())
}

// GetDateString : 입력 받은 구분자를 사용하여 만든 날짜 문자열 리턴
func GetDateString(date time.Time, sep string) string {
	return fmt.Sprintf("%02d%s%02d%s%02d",
		date.Year(), sep, date.Month(), sep, date.Day())
}

// GetYearMonthString : 입력 받은 구분자를 사용하여 년,월 문자열 리턴
func GetYearMonthString(date time.Time, sep string) string {
	return fmt.Sprintf("%02d%s%02d", date.Year(), sep, date.Month())
}

// GetYearsInRange : 시작일과 종료일 사이의 년도 획득(시작, 종료년도 포함)
func GetYearsInRange(start time.Time, end time.Time) []int {
	var result []int
	startYear := start.Year()
	endYear := end.Year()

	for year := startYear; year <= endYear; year++ {
		result = append(result, year)
	}

	return result
}
