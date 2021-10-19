package util

import (
	"bufio"
	"encoding/csv"
	"log"
	"media_crawling/alarm"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// MakeDirIfNotExists : 디렉토리가 없다면 생성
func MakeDirIfNotExists(dirPath string) {
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dirPath, os.ModePerm)
		}
	}
}

// MakePath : 문자열 배열을 받아 조합하여 경로 생성
func MakePath(pathElems []string) string {
	path := strings.Join(pathElems, "/")

	return filepath.FromSlash(path)
}

// GetWorkingDirPath : Working Directory 경로 얻기
func GetWorkingDirPath() string {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(filepath.Dir(b))
	)

	return basepath
}

// ReadCsv : csv 파일 읽기
func ReadCsv(path string) [][]string {
	// 파일 오픈
	file, _ := os.Open(path)
	defer file.Close()

	// csv reader 생성
	rdr := csv.NewReader(bufio.NewReader(file))

	// csv 내용 모두 읽기
	rows, _ := rdr.ReadAll()

	return rows
}

// WriteCsv : csv 파일 쓰기
func WriteCsv(path string, rows [][]string) {
	// 파일 생성
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	// csv writer 생성
	wr := csv.NewWriter(bufio.NewWriter(file))

	// csv 내용 쓰기
	for i := range rows {
		wr.Write(rows[i])
	}
	wr.Flush()
}

// GetFileScanner : bufio.Scanner를 획득
func GetFileScanner(filePath string) *bufio.Scanner {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	return scanner
}
