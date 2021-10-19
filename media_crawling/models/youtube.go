package models

import (
	"github.com/fatih/structs"
)

// Youtube : 유튜브 비디오 그조체
type Youtube struct {
	Query        string `gorm:"column:Query"`
	VideoID      string `gorm:"column:VideoID"`
	ChannelID    string `gorm:"column:ChannelID"`
	ChannelTitle string `gorm:"column:ChannelTitle"`
	Title        string `gorm:"column:Title"`
	Description  string `gorm:"column:Description"`
	ThumbnailURL string `gorm:"column:ThumbnailURL"`
	PublishedAt  string `gorm:"column:PublishedAt"`
}

// TableName : 언론사 테이블 명
func (Youtube) TableName() string {
	return "AGRI_YOUTUBE"
}

// Columns : 언론사 컬럼 명
func (Youtube) Columns() []string {
	return structs.Names(&Youtube{})
}
