package models

import "github.com/fatih/structs"

// Blog : 블로그 구조체
type Blog struct {
	Query       string `gorm:"column:Query"`
	Title       string `gorm:"column:Title"`
	Link        string `gorm:"column:Link"`
	Bloggerlink string `gorm:"column:Bloggerlink"`
	Description string `gorm:"column:Description"`
	Bloggername string `gorm:"column:Bloggername"`
	Postdate    string `gorm:"column:Postdate"`
}

// TableName : 블로그 테이블 명
func (Blog) TableName() string {
	return "AGRI_BLOG"
}

// Columns : 블로그 컬럼 명
func (Blog) Columns() []string {
	return structs.Names(&Blog{})
}
