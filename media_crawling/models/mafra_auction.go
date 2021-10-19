package models

import "github.com/fatih/structs"

// MafraAdjAuctionStats : mafra 일별 정산 경락 요약정보 구조체
type MafraAdjAuctionStats struct {
	AUCNGDE           string `gorm:"column:AucngDe"`
	PBLMNGWHSALMRKTNM string `gorm:"column:PblmngWhsalMrktNm"`
	PBLMNGWHSALMRKTCD string `gorm:"column:PblmngWhsalMrktCd"`
	CPRNM             string `gorm:"column:CprNm"`
	CPRCD             string `gorm:"column:CprCd"`
	PRDLSTNM          string `gorm:"column:PrdlstNm"`
	PRDLSTCD          string `gorm:"column:PrdlstCd"`
	SPCIESNM          string `gorm:"column:SpciesNm"`
	SPCIESCD          string `gorm:"column:SpciesCd"`
	DELNGBUNDLEQY     string `gorm:"column:DelngbundleQy"`
	STNDRD            string `gorm:"column:Stndrd"`
	STNDRDCD          string `gorm:"column:StndrdCd"`
	GRAD              string `gorm:"column:Grad"`
	GRADCD            string `gorm:"column:GradCd"`
	SANJICD           string `gorm:"column:SanjiCd"`
	SANJINM           string `gorm:"column:SanjiNm"`
	MUMMAMT           string `gorm:"column:MummAmt"`
	AVRGAMT           string `gorm:"column:AvrgAmt"`
	MXMMAMT           string `gorm:"column:MxmmAmt"`
	DELNGQY           string `gorm:"column:DelngQy"`
	CNTS              string `gorm:"column:Cnts"`
}

// MafraAdjAuctionQuantity : mafra 일별 정산 경락 거래량 구조체
type MafraAdjAuctionQuantity struct {
	AuctionDate    string `gorm:"column:AuctionDate"`
	StdItemName    string `gorm:"column:StdItemName"`
	StdItemCode    string `gorm:"column:StdItemCode"`
	StdSpeciesName string `gorm:"column:StdSpeciesName"`
	StdSpeciesCode string `gorm:"column:StdSpeciesCode"`
	Quantity       string `gorm:"column:Quantity"`
}

// TableName :mafra 일별 정산 경락 요약정보 테이블 명
func (MafraAdjAuctionStats) TableName() string {
	return "MAFRA_ADJ_AUCTION_STATS"
}

// Columns :mafra 일별 정산 경락 요약정보 컬럼 명
func (MafraAdjAuctionStats) Columns() []string {
	return structs.Names(&MafraAdjAuctionStats{})
}

// TableName :mafra 일별 정산 경락 거래량 테이블 명
func (MafraAdjAuctionQuantity) TableName() string {
	return "MAFRA_ADJ_AUCTION_QUANTITY"
}

// Columns :mafra 일별 정산 경락 거래량 컬럼 명
func (MafraAdjAuctionQuantity) Columns() []string {
	return structs.Names(&MafraAdjAuctionQuantity{})
}
