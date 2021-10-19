package main

import (
	"media_crawling/checker"
	"media_crawling/config"
	"media_crawling/crawler"
	"media_crawling/util"
	"os"
)

func main() {
	conf := config.Conf
	logFilePath := util.MakePath([]string{config.Conf.SubLogDirPath, "youtubeByChannel.txt"})
	logErrorMessage := "채널 별 유튜브 수집 문제 생김"

	os.Remove(logFilePath)
	util.MakeDirIfNotExists(conf.LogDirPath)
	util.MakeDirIfNotExists(conf.SubLogDirPath)

	crawler.CrawlYoutubeByChannel(logFilePath)
	checker.CheckLog(logFilePath, logErrorMessage)
}
