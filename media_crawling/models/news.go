package models

import "github.com/fatih/structs"

// News : 뉴스기사 구조체
type News struct {
	Query        string `gorm:"column:Query"`
	Title        string `gorm:"column:Title"`
	Link         string `gorm:"column:Link"`
	PressKeyword string `gorm:"column:PressKeyword"`
	Description  string `gorm:"column:Description"`
	PubDate      string `gorm:"column:PubDate"`
}

// TableName : 뉴스기사 테이블 명
func (News) TableName() string {
	return "AGRI_NEWS"
}

// Columns : 뉴스 기사 컬럼 명
func (News) Columns() []string {
	return structs.Names(&News{})
}
