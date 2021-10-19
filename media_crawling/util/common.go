package util

import (
	"errors"
	"fmt"
	"media_crawling/alarm"
	"strings"
)

// ReplaceString : uselessWord를 지정된 문자열로 변경
func ReplaceString(source string, uselessWords map[string]string) string {
	result := source
	for uselessWord, replacement := range uselessWords {
		result = strings.Replace(result, uselessWord, replacement, -1)
	}

	return result
}

// InArray : 입력 받은 인자가 Array 안에 있는지 체크
func InArray(target string, array []string) bool {
	for _, item := range array {
		if item == target {
			return true
		}
	}

	return false
}

// IsLastElement : 마지막 인자인지 체크
func IsLastElement(index int, length int) bool {
	if index == length-1 {
		return true
	}

	return false
}

// AllTrue : inputArray가 모두 True 인지 확인
func AllTrue(inputArray []bool) bool {
	for _, elem := range inputArray {
		if !elem {
			return false
		}
	}

	return true
}

// GetMaxPageNo : 최대 페이지 수 획득
func GetMaxPageNo(numerator int, denominator int) int {
	quotient := numerator / denominator
	remainder := numerator % denominator
	maxPageNo := quotient
	if remainder != 0 {
		maxPageNo++
	}

	return maxPageNo
}

// CheckError : Error를 받아 nil이 아니면 slack으로 오류메세지 전송 후 panic
func CheckError(err error, functionName string) {
	if err != nil {
		errorMessage := fmt.Sprintf("%s :: 에러 : %s", functionName, err.Error())
		alarm.PostMessage("default", errorMessage)
		panic(err)
	}
}

// CheckCondition : condition을 받아 false 면 slack으로 오류메세지 전송 후 panic
func CheckCondition(condition bool, functionName string, message string) {
	if condition {
		errorMessage := fmt.Sprintf("%s :: 에러 : %s", functionName, message)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}
}
