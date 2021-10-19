package models

import "github.com/fatih/structs"

// Cafe : 카페 구조체
type Cafe struct {
	Query       string `gorm:"column:Query"`
	Title       string `gorm:"column:Title"`
	Link        string `gorm:"column:Link"`
	Cafeurl     string `gorm:"column:Cafeurl"`
	Description string `gorm:"column:Description"`
	Cafename    string `gorm:"column:Cafename"`
}

// TableName :카페 테이블 명
func (Cafe) TableName() string {
	return "AGRI_CAFE"
}

// Columns :카페 컬럼 명
func (Cafe) Columns() []string {
	return structs.Names(&Cafe{})
}
