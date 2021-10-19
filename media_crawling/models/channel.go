package models

import "github.com/fatih/structs"

// Channel : 채널 구조체
type Channel struct {
	ID    string `gorm:"column:ID"`
	Title string `gorm:"column:Title"`
}

// TableName : 채널 테이블 명
func (Channel) TableName() string {
	return "AGRI_CHANNEL"
}

// Columns : 채널 컬럼 명
func (Channel) Columns() []string {
	return structs.Names(&Channel{})
}
