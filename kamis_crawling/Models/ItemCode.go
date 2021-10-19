package models

import "kamis_crawling/common"

// Pcode 코드
type Pcode struct {
	ItemCode     int    `gorm:"column:ItemCode"`
	ItemKindCode string `gorm:"column:ItemKindCode"`
	GradeCode    string `gorm:"column:GradeCode"`
}

// GetItemCodeList 품목 코드 조회
func GetItemCodeList(cls string) []Pcode {
	db := common.GetDB()

	var list []Pcode

	if cls == "01" {
		db.Raw(`SELECT ItemCode, ItemKindCode, GradeCode
		from AGRI_ITEM_MNG
		where ShipType = 'retail'
		group by ItemCode, ItemKindCode, GradeCode`).Scan(&list)
	} else if cls == "02" {
		db.Raw(`SELECT ItemCode, ItemKindCode, GradeCode
		from AGRI_ITEM_MNG
		where ShipType = 'whole'
		group by ItemCode, ItemKindCode, GradeCode`).Scan(&list)
	}

	return list
}
