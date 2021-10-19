package models

import "github.com/fatih/structs"

// RawBreakingAuction : 실시간 경락 속보 원본 구조체
type RawBreakingAuction struct {
	Bidtime    string `gorm:"column:Bidtime"`
	Coname     string `gorm:"column:Coname"`
	Gradename  string `gorm:"column:Gradename"`
	Marketname string `gorm:"column:Marketname"`
	Mclassname string `gorm:"column:Mclassname"`
	Price      string `gorm:"column:Price"`
	Sanji      string `gorm:"column:Sanji"`
	Sclassname string `gorm:"column:Sclassname"`
	Tradeamt   string `gorm:"column:Tradeamt"`
	Unitname   string `gorm:"column:Unitname"`
	Chulagtnm  string `gorm:"column:Chulagtnm"`
}

// BreakingAuction : 실시간 경락 속보 구조체
type BreakingAuction struct {
	Bidtime    string `gorm:"column:Bidtime"`
	Coname     string `gorm:"column:Coname"`
	Gradename  string `gorm:"column:Gradename"`
	Marketname string `gorm:"column:Marketname"`
	Mclassname string `gorm:"column:Mclassname"`
	Price      string `gorm:"column:Price"`
	Sanji      string `gorm:"column:Sanji"`
	Sclassname string `gorm:"column:Sclassname"`
	Tradeamt   string `gorm:"column:Tradeamt"`
	Unitname   string `gorm:"column:Unitname"`
	Chulagtnm  string `gorm:"column:Chulagtnm"`
	Cocode     string `gorm:"column:Cocode"`
	Gradecode  string `gorm:"column:Gradecode"`
	Marketco   string `gorm:"column:Marketco"`
	MclassCode string `gorm:"column:MclassCode"`
	Zipcode    string `gorm:"column:Zipcode"`
	SClassCode string `gorm:"column:SClassCode"`
	Unitamt    string `gorm:"column:Unitamt"`
	Unitcode   string `gorm:"column:Unitcode"`
}

// TableName : 실시간 경락 속보 원본 테이블 명
func (RawBreakingAuction) TableName() string {
	return "RAW_BREAKING_AUCTION"
}

// Columns : 실시간 경락 속보 원본 컬럼 명
func (RawBreakingAuction) Columns() []string {
	return structs.Names(&RawBreakingAuction{})
}

// TableName : 실시간 경락 속보 테이블 명
func (BreakingAuction) TableName() string {
	return "BREAKING_AUCTION"
}

// Columns : 실시간 경락 속보 컬럼 명
func (BreakingAuction) Columns() []string {
	return structs.Names(&BreakingAuction{})
}
