package models

import "github.com/fatih/structs"

// Weather : 일별 날씨 구조체
type Weather struct {
	AreaID        string `gorm:"column:AreaId"`
	AreaName      string `gorm:"column:AreaName"`
	DayAvgRhm     string `gorm:"column:DayAvgRhm"`
	DayAvgTa      string `gorm:"column:DayAvgTa"`
	DayAvgWs      string `gorm:"column:DayAvgWs"`
	DayMaxTa      string `gorm:"column:DayMaxTa"`
	DayMinRhm     string `gorm:"column:DayMinRhm"`
	DayMinTa      string `gorm:"column:DayMinTa"`
	DaySumRn      string `gorm:"column:DaySumRn"`
	DaySumSs      string `gorm:"column:DaySumSs"`
	PaCropName    string `gorm:"column:PaCropName"`
	PaCropSpeID   string `gorm:"column:PaCropSpeId"`
	PaCropSpeName string `gorm:"column:PaCropSpeName"`
	WmCd          string `gorm:"column:WmCd"`
	WmCount       string `gorm:"column:WmCount"`
	Ymd           string `gorm:"column:Ymd"`
}

// TableName : 일별 날씨 테이블 명
func (Weather) TableName() string {
	return "AGRI_WEATHER"
}

// Columns : 일별 날씨 컬럼 명
func (Weather) Columns() []string {
	return structs.Names(&Weather{})
}
