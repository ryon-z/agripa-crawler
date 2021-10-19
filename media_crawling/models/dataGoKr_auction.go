package models

import "github.com/fatih/structs"

// RealtimeAuction : 실시간 경락 도매시장 유통통계
type RealtimeAuction struct {
	CprCode          string `gorm:"column:CprCode"`
	CprInsttCode     string `gorm:"column:CprInsttCode"`
	CprInsttNm       string `gorm:"column:CprInsttNm"`
	CprNm            string `gorm:"column:CprNm"`
	DelngDe          string `gorm:"column:DelngDe"`
	DelngPric        string `gorm:"column:DelngPric"`
	DelngPrutCode    string `gorm:"column:DelngPrutCode"`
	DelngQy          string `gorm:"column:DelngQy"`
	LclasCode        string `gorm:"column:LclasCode"`
	LclasNm          string `gorm:"column:LclasNm"`
	MlsfcCode        string `gorm:"column:MlsfcCode"`
	MlsfcNm          string `gorm:"column:MlsfcNm"`
	SbidPricAvg      string `gorm:"column:SbidPricAvg"`
	SbidPricMax      string `gorm:"column:SbidPricMax"`
	SbidPricMin      string `gorm:"column:SbidPricMin"`
	SbidPricMvAvg    string `gorm:"column:SbidPricMvAvg"`
	StdFrmlcNewCode  string `gorm:"column:StdFrmlcNewCode"`
	StdFrmlcNewNm    string `gorm:"column:StdFrmlcNewNm"`
	StdGradCode      string `gorm:"column:StdGradCode"`
	StdGradNm        string `gorm:"column:StdGradNm"`
	StdMgNewCode     string `gorm:"column:StdMgNewCode"`
	StdMgNewNm       string `gorm:"column:StdMgNewNm"`
	StdMtcCode       string `gorm:"column:StdMtcCode"`
	StdMtcNewCode    string `gorm:"column:StdMtcNewCode"`
	StdMtcNewNm      string `gorm:"column:StdMtcNewNm"`
	StdMtcNm         string `gorm:"column:StdMtcNm"`
	StdPrdlstCode    string `gorm:"column:StdPrdlstCode"`
	StdPrdlstNewCode string `gorm:"column:StdPrdlstNewCode"`
	StdPrdlstNewNm   string `gorm:"column:StdPrdlstNewNm"`
	StdPrdlstNm      string `gorm:"column:StdPrdlstNm"`
	StdQlityNewCode  string `gorm:"column:StdQlityNewCode"`
	StdQlityNewNm    string `gorm:"column:StdQlityNewNm"`
	StdUnitCode      string `gorm:"column:StdUnitCode"`
	StdUnitNewCode   string `gorm:"column:StdUnitNewCode"`
	StdUnitNewNm     string `gorm:"column:StdUnitNewNm"`
	StdUnitNm        string `gorm:"column:StdUnitNm"`
	WhsalCode        string `gorm:"column:WhsalCode"`
	WhsalMrktCode    string `gorm:"column:WhsalMrktCode"`
	WhsalMrktNm      string `gorm:"column:WhsalMrktNm"`
	WhsalNm          string `gorm:"column:WhsalNm"`
}

// AdjustedAuction : 정산 경락 도매시장 유통통계
type AdjustedAuction struct {
	CprCode          string `gorm:"column:CprCode"`
	CprInsttCode     string `gorm:"column:CprInsttCode"`
	CprInsttNm       string `gorm:"column:CprInsttNm"`
	CprNm            string `gorm:"column:CprNm"`
	DelngDe          string `gorm:"column:DelngDe"`
	DelngPric        string `gorm:"column:DelngPric"`
	DelngPrutCode    string `gorm:"column:DelngPrutCode"`
	DelngQy          string `gorm:"column:DelngQy"`
	LclasCode        string `gorm:"column:LclasCode"`
	LclasNm          string `gorm:"column:LclasNm"`
	MlsfcCode        string `gorm:"column:MlsfcCode"`
	MlsfcNm          string `gorm:"column:MlsfcNm"`
	SbidPricAvg      string `gorm:"column:SbidPricAvg"`
	SbidPricMax      string `gorm:"column:SbidPricMax"`
	SbidPricMin      string `gorm:"column:SbidPricMin"`
	SbidPricMvAvg    string `gorm:"column:SbidPricMvAvg"`
	StdFrmlcNewCode  string `gorm:"column:StdFrmlcNewCode"`
	StdFrmlcNewNm    string `gorm:"column:StdFrmlcNewNm"`
	StdGradCode      string `gorm:"column:StdGradCode"`
	StdGradNm        string `gorm:"column:StdGradNm"`
	StdMgNewCode     string `gorm:"column:StdMgNewCode"`
	StdMgNewNm       string `gorm:"column:StdMgNewNm"`
	StdMtcCode       string `gorm:"column:StdMtcCode"`
	StdMtcNewCode    string `gorm:"column:StdMtcNewCode"`
	StdMtcNewNm      string `gorm:"column:StdMtcNewNm"`
	StdMtcNm         string `gorm:"column:StdMtcNm"`
	StdPrdlstCode    string `gorm:"column:StdPrdlstCode"`
	StdPrdlstNewCode string `gorm:"column:StdPrdlstNewCode"`
	StdPrdlstNewNm   string `gorm:"column:StdPrdlstNewNm"`
	StdPrdlstNm      string `gorm:"column:StdPrdlstNm"`
	StdQlityNewCode  string `gorm:"column:StdQlityNewCode"`
	StdQlityNewNm    string `gorm:"column:StdQlityNewNm"`
	StdUnitCode      string `gorm:"column:StdUnitCode"`
	StdUnitNewCode   string `gorm:"column:StdUnitNewCode"`
	StdUnitNewNm     string `gorm:"column:StdUnitNewNm"`
	StdUnitNm        string `gorm:"column:StdUnitNm"`
	WhsalCode        string `gorm:"column:WhsalCode"`
	WhsalMrktCode    string `gorm:"column:WhsalMrktCode"`
	WhsalMrktNm      string `gorm:"column:WhsalMrktNm"`
	WhsalNm          string `gorm:"column:WhsalNm"`
}

// TableName : 정산 경락 도매시장 유통통계 테이블 명
func (AdjustedAuction) TableName() string {
	return "AGRI_ADJ_AUCTION"
}

// Columns : 정산 경락 도매시장 유통통계 컬럼 명
func (AdjustedAuction) Columns() []string {
	return structs.Names(&AdjustedAuction{})
}

// TableName : 실시간 경락 도매시장 유통통계 테이블 명
func (RealtimeAuction) TableName() string {
	return "AGRI_REAL_AUCTION"
}

// Columns : 실시간 경락 도매시장 유통통계 컬럼 명
func (RealtimeAuction) Columns() []string {
	return structs.Names(&RealtimeAuction{})
}
