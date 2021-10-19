package main

import (
	"errors"
	"fmt"
	"media_crawling/alarm"
	"media_crawling/checker"
	"media_crawling/config"
	"media_crawling/crawler"
	"media_crawling/util"
	"os"
	"strconv"
	"strings"
)

func main() {
	conf := config.Conf
	logFilePath := util.MakePath([]string{config.Conf.SubLogDirPath, "youtubeByQuery.txt"})
	logErrorMessage := "쿼리 별 유튜브 수집 문제 생김"

	os.Remove(logFilePath)
	util.MakeDirIfNotExists(conf.LogDirPath)
	util.MakeDirIfNotExists(conf.SubLogDirPath)

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 1 {
		message := fmt.Sprintf("인자가 1개가 아닙니다. 인자 : %s", strings.Join(argsWithoutProg, ", "))
		err := errors.New(message)
		alarm.PostMessage("default", message)
		panic(err)
	}

	firstCodeNum, err := strconv.Atoi(argsWithoutProg[0])
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	crawler.CrawlYoutube(firstCodeNum, logFilePath)
	checker.CheckLog(logFilePath, logErrorMessage)
}
