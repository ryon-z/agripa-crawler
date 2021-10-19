package main

import (
	"fmt"
	"media_crawling/config"
	"media_crawling/refiner"
	"time"
)

func main() {
	conf := config.Conf
	loc := conf.Timezone

	todayTime := time.Now().In(loc)
	numAllowedDays := 3
	datetimeToLeave := todayTime.AddDate(0, 0, -1*numAllowedDays)
	dateString := fmt.Sprintf(
		"%02d-%02d-%02d",
		datetimeToLeave.Year(),
		datetimeToLeave.Month(),
		datetimeToLeave.Day())

	refiner.RemoveOldBreakingAuction(dateString)
}
