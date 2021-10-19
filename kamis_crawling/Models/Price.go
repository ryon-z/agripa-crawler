package models

import (
	"kamis_crawling/common"
)

// WholePrice 도매 가격 정보
type WholePrice struct {
	ShipDate     string  `gorm:"column:ShipDate"`
	ItemCode     int     `gorm:"column:ItemCode"`
	ItemKindCode string  `gorm:"column:ItemKindCode"`
	GradeCode    string  `gorm:"column:GradeCode"`
	MarketName   string  `gorm:"column:MarketName"`
	AreaName     string  `gorm:"column:AreaName"`
	ShipPrice    int     `gorm:"column:ShipPrice"`
	ShipUnit     string  `gorm:"column:ShipUnit"`
	ShipAmt      float64 `gorm:"column:ShipAmt"`
}

// TableName 도매 가격 정보 테이블명
func (WholePrice) TableName() string {
	return "AGRI_WHOLE_PRICE"
}

// RetailPrice 소매 가격 정보
type RetailPrice struct {
	ShipDate     string  `gorm:"column:ShipDate"`
	ItemCode     int     `gorm:"column:ItemCode"`
	ItemKindCode string  `gorm:"column:ItemKindCode"`
	GradeCode    string  `gorm:"column:GradeCode"`
	MarketName   string  `gorm:"column:MarketName"`
	AreaName     string  `gorm:"column:AreaName"`
	ShipPrice    int     `gorm:"column:ShipPrice"`
	ShipUnit     string  `gorm:"column:ShipUnit"`
	ShipAmt      float64 `gorm:"column:ShipAmt"`
}

// TableName 소매 가격 정보 테이블명
func (RetailPrice) TableName() string {
	return "AGRI_RETAIL_PRICE"
}

// AddWholePrice 도매가격 입력
func AddWholePrice(row WholePrice) int {
	db := common.GetDB()

	db.Create(row)

	return 0
}

// AddRetailPrice 도매가격 입력
func AddRetailPrice(row RetailPrice) int {
	db := common.GetDB()

	db.Create(row)

	return 0
}
