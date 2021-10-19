package util_test

import (
	"fmt"
	"media_crawling/util"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type fileTestData struct {
	path string
	data [][]string
}

func (d *fileTestData) setData() {
	pwd := util.GetWorkingDirPath()
	path := pwd + "/test.csv"
	d.path = filepath.FromSlash(path)
	d.data = [][]string{
		{"a", "b", "c"},
		{"d", "e", "f"},
	}
}

func clean() {
	var testData fileTestData
	testData.setData()

	err := os.Remove(testData.path)
	if err != nil {
		panic(err)
	}

	fmt.Println("clean")
}

func TestWriteCsv(t *testing.T) {
	var testData fileTestData
	testData.setData()

	util.WriteCsv(testData.path, testData.data)
	if _, err := os.Stat(testData.path); os.IsNotExist(err) {
		t.Error("csv 쓰기 실패")
	}
	fmt.Println("TestWriteCsv")
}

func TestReadCsv(t *testing.T) {
	var testData fileTestData
	testData.setData()

	loadedCsv := util.ReadCsv(testData.path)
	if !reflect.DeepEqual(testData.data, loadedCsv) {
		t.Error("csv 파일 읽기 실패")
	}
	fmt.Println("TestReadCsv")
}

func TestMain(m *testing.M) {
	exitVal := m.Run()
	clean()
	os.Exit(exitVal)
}
