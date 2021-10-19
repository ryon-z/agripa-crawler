package checker

import (
	"media_crawling/alarm"
	"media_crawling/config"
	"media_crawling/util"
	"strings"

	"github.com/slack-go/slack"
)

type checkerMeta struct {
	start string
	end   string
}

// isLogNormal : 로그 파일에서 "시작"과 "끝" 수가 서로 맞는지 확인
func isLogNormal(logFilePath string, meta checkerMeta) bool {
	scanner := util.GetFileScanner(logFilePath)
	startPhrase := meta.start
	endPhrase := meta.end
	startCount := 0
	endCount := 0

	var row string
	for scanner.Scan() {
		row = scanner.Text()
		startCount += isContained(row, startPhrase)
		endCount += isContained(row, endPhrase)
	}

	return startCount == endCount
}

func isContained(s string, substr string) int {
	if strings.Contains(s, substr) {
		return 1
	}

	return 0
}

// CheckLog : 로그가 이상이 없는지 체크
func CheckLog(logFilePath string, errorMessage string) {
	checkerMetas := []checkerMeta{
		{"수집 시작", "수집 종료"},
		{"정상시작", "정상종료"},
	}

	for _, checkerMeta := range checkerMetas {
		if !isLogNormal(logFilePath, checkerMeta) {
			alarm.SlackBot.PostMessage(
				config.Secret["slack:channel_id"],
				slack.MsgOptionText(errorMessage, true),
			)
			return
		}
	}
}
