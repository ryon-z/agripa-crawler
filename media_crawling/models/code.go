package models

import (
	"github.com/fatih/structs"
)

// GarakCode : 가락시장 코드 구조체
type GarakCode struct {
	GarrakCode string `gorm:"column:GarrakCode"`
	GarrakName string `gorm:"column:GarrakName"`
	Sclasscode string `gorm:"column:Sclasscode"`
	StanCode   string `gorm:"column:StanCode"`
}

// WholesaleMarketCode : 도매시장 코드 구조체
type WholesaleMarketCode struct {
	Marketco string `gorm:"column:MarketCo"`
	Marketnm string `gorm:"column:MarketNm"`
}

// WholesaleMarketCoCode : 도매시장 법인 코드 구조체
type WholesaleMarketCoCode struct {
	Cocode     string `gorm:"column:CoCode"`
	Coname     string `gorm:"column:CoName"`
	Marketcode string `gorm:"column:MarketCode"`
	Marketname string `gorm:"column:MarketName"`
}

// StdGradeCode : 표준 등급 코드 구조체
type StdGradeCode struct {
	Gradecode string `gorm:"column:GradeCode"`
	Gradename string `gorm:"column:GradeName"`
}

// StdUnitCode : 표준 단위 코드 구조체
type StdUnitCode struct {
	Unitcode string `gorm:"column:UnitCode"`
	Unitname string `gorm:"column:UnitName"`
}

// PlaceOriginCode : 산지 코드 구조체
type PlaceOriginCode struct {
	Zipcode string `gorm:"column:Zipcode"`
	Sido    string `gorm:"column:Sido"`
	Sigun   string `gorm:"column:Sigun"`
	Dong    string `gorm:"column:Dong"`
}

// StdSpeciesCode : 표준품종코드 구조체
type StdSpeciesCode struct {
	LClassCode string `gorm:"column:LClassCode"`
	LClassName string `gorm:"column:LClassName"`
	MClassCode string `gorm:"column:MClassCode"`
	MClassName string `gorm:"column:MClassName"`
	SClassCode string `gorm:"column:SClassCode"`
	SClassName string `gorm:"column:SClassName"`
}

// StdItemCode : 표준품목코드 구조체
type StdItemCode struct {
	ItemCode string `gorm:"column:ItemCode"`
	ItemName string `gorm:"column:ItemName"`
}

// ItemMapping : 품목 매핑 테이블 구조체
type ItemMapping struct {
	StdItemCode    string `gorm:"column:StdItemCode"`
	ExaminItemCode string `gorm:"column:ExaminItemCode"`
	HskPrdlstCode  string `gorm:"column:HskPrdlstCode"`
}

// StdItemKeyword : 표준품목별 검색어 구조체
type StdItemKeyword struct {
	ItemCode       string `gorm:"column:ItemCode"`
	ExposedKeyword string `gorm:"column:ExposedKeyword"`
	Keyword        string `gorm:"column:Keyword"`
	Priority       string `gorm:"column:Priority"`
	IsDisplay      string `gorm:"column:IsDisplay"`
	NumSearch      string `gorm:"column:NumSearch"`
}

// TableName : 가락시장 코드 테이블 명
func (GarakCode) TableName() string {
	return "GARAK_CODE"
}

// Columns : 가락시장 코드 컬럼 명
func (GarakCode) Columns() []string {
	return structs.Names(&GarakCode{})
}

// TableName : 도매시장 코드 테이블 명
func (WholesaleMarketCode) TableName() string {
	return "WHOLESALE_MARKET_CODE"
}

// Columns : 도매시장 코드 컬럼 명
func (WholesaleMarketCode) Columns() []string {
	return structs.Names(&WholesaleMarketCode{})
}

// TableName : 도매시장 법인 코드 테이블 명
func (WholesaleMarketCoCode) TableName() string {
	return "WHOLESALE_MARKET_CO_CODE"
}

// Columns : 도매시장 법인 코드 컬럼 명
func (WholesaleMarketCoCode) Columns() []string {
	return structs.Names(&WholesaleMarketCoCode{})
}

// TableName : 표준 등급 코드 테이블 명
func (StdGradeCode) TableName() string {
	return "STD_GRADE_CODE"
}

// Columns : 표준 등급 코드 컬럼 명
func (StdGradeCode) Columns() []string {
	return structs.Names(&StdGradeCode{})
}

// TableName : 표준 단위 코드 테이블 명
func (StdUnitCode) TableName() string {
	return "STD_UNIT_CODE"
}

// Columns : 표준 단위 코드 컬럼 명
func (StdUnitCode) Columns() []string {
	return structs.Names(&StdUnitCode{})
}

// TableName : 산지 코드 테이블 명
func (PlaceOriginCode) TableName() string {
	return "PLACE_ORIGIN_CODE"
}

// Columns : 산지 코드 컬럼 명
func (PlaceOriginCode) Columns() []string {
	return structs.Names(&PlaceOriginCode{})
}

// TableName : 표준품종코드 테이블 명
func (StdSpeciesCode) TableName() string {
	return "STD_SPECIES_CODE"
}

// Columns : 표준품종코드 컬럼 명
func (StdSpeciesCode) Columns() []string {
	return structs.Names(&StdSpeciesCode{})
}

// TableName : 표준품목코드 테이블 명
func (StdItemCode) TableName() string {
	return "STD_ITEM_CODE"
}

// Columns : 표준품목코드 컬럼 명
func (StdItemCode) Columns() []string {
	return structs.Names(&StdItemCode{})
}

// TableName : 품목 매핑 테이블 테이블 명
func (ItemMapping) TableName() string {
	return "ITEM_MAPPING"
}

// Columns : 품목 매핑 테이블 컬럼 명
func (ItemMapping) Columns() []string {
	return structs.Names(&ItemMapping{})
}

// TableName : 표준품목별 검색어 테이블 명
func (StdItemKeyword) TableName() string {
	return "STD_ITEM_KEYWORD"
}

// Columns : 표준품목별 검색어 컬럼 명
func (StdItemKeyword) Columns() []string {
	return structs.Names(&StdItemKeyword{})
}
