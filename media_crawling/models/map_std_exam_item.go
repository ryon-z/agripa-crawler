package models

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

// MapStdExamItem : 표준코드-조사코드 매핑 구조체
type MapStdExamItem struct {
	CatgoryNewCode   string `gorm:"column:CatgoryNewCode"`
	CatgoryNewNm     string `gorm:"column:CatgoryNewNm"`
	StdPrdlstCode    string `gorm:"column:StdPrdlstCode"`
	StdPrdlstNm      string `gorm:"column:StdPrdlstNm"`
	ExaminPrdlstCode string `gorm:"column:ExaminPrdlstCode"`
	ExaminPrdlstNm   string `gorm:"column:ExaminPrdlstNm"`
	StdSpciesCode    string `gorm:"column:StdSpciesCode"`
	StdSpciesNm      string `gorm:"column:StdSpciesNm"`
	ExaminSpciesCode string `gorm:"column:ExaminSpciesCode"`
	ExaminSpciesNm   string `gorm:"column:ExaminSpciesNm"`
	UpdtDe           string `gorm:"column:UpdtDe"`
}

// TableName : 표준코드-조사코드 매핑 테이블 명
func (MapStdExamItem) TableName() string {
	return "MAP_STD_EXAM_ITEM"
}

// Columns : 표준코드-조사코드 매핑 컬럼 명
func (MapStdExamItem) Columns() []string {
	return structs.Names(&MapStdExamItem{})
}

// DistinctStdSpciesCodeResponse : 조사코드와 매핑된 표준품종코드 응답결과 구조체
type DistinctStdSpciesCodeResponse struct {
	StdSpciesCode string `gorm:"column:StdSpciesCode"`
}

// GetStdSpciesCodesMappedExam : 중복이 제거된 조사코드와 매핑된 표준품종코드 목록 획득
func GetStdSpciesCodesMappedExam(db *gorm.DB) []DistinctStdSpciesCodeResponse {
	var maps []DistinctStdSpciesCodeResponse
	tableName := MapStdExamItem{}.TableName()
	sqlQuery := fmt.Sprintf("SELECT distinct StdSpciesCode FROM %s", tableName)
	db.Raw(sqlQuery).Scan(&maps)

	return maps
}
