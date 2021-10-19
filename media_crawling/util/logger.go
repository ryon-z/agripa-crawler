package util

import (
	"log"
	"media_crawling/alarm"
	"os"
)

// GetFileLogger : 파일 로거를 리턴
func GetFileLogger(logPath string) (*os.File, *log.Logger) {
	fpLog, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	logger := log.New(fpLog, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	return fpLog, logger
}

// CloseFileLogger : 파일 로거를 종료
func CloseFileLogger(fpLog *os.File) {
	fpLog.Close()
}
