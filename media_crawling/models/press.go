package models

import (
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

// Press : 언론사 구조체
type Press struct {
	Name    string `gorm:"column:Name"`
	Link    string `gorm:"column:Link"`
	Keyword string `gorm:"column:Keyword"`
}

// TableName : 언론사 테이블 명
func (Press) TableName() string {
	return "AGRI_PRESS"
}

// Columns : 언론사 컬럼 명
func (Press) Columns() []string {
	return structs.Names(&Press{})
}

// GetPresses : Press 테이블 rows 획득
func GetPresses(db *gorm.DB) []Press {
	var presses []Press
	db.Find(&presses)

	return presses
}
