package main

import (
	"media_crawling/config"
	"media_crawling/crawler"
	"media_crawling/util"
	"time"
)

func main() {
	conf := config.Conf
	loc := conf.Timezone

	todayTime := time.Now().In(loc)
	todayString := util.GetDateStringNoSep(todayTime) // yyyymmdd

	crawler.CrawlBreakingAuction(todayString)
}
