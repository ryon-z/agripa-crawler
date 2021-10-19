package models

import "github.com/fatih/structs"

// MafraExamination : mafra 농수축산 유통정보 조사가격(농수축산물표준코드변환)
type MafraExamination struct {
	EXAMINDE         string `gorm:"column:ExaminDe"`
	EXAMINSENM       string `gorm:"column:ExaminSeNm"`
	EXAMINSECODE     string `gorm:"column:ExaminSeCode"`
	EXAMINAREANAME   string `gorm:"column:ExaminAreaName"`
	EXAMINAREACODE   string `gorm:"column:ExaminAreaCode"`
	EXAMINMRKTNM     string `gorm:"column:ExaminMrktNm"`
	EXAMINMRKTCODE   string `gorm:"column:ExaminMrktCode"`
	STDMRKTNM        string `gorm:"column:StdMrktNm"`
	STDMRKTCODE      string `gorm:"column:StdMrktCode"`
	EXAMINPRDLSTNM   string `gorm:"column:ExaminPrdlstNm"`
	EXAMINPRDLSTCODE string `gorm:"column:ExaminPrdlstCode"`
	EXAMINSPCIESNM   string `gorm:"column:ExaminSpciesNm"`
	EXAMINSPCIESCODE string `gorm:"column:ExaminSpciesCode"`
	STDLCLASNM       string `gorm:"column:StdLclasNm"`
	STDLCLASCO       string `gorm:"column:StdLclasco"`
	STDPRDLSTNM      string `gorm:"column:StdPrdlstNm"`
	STDPRDLSTCODE    string `gorm:"column:StdPrdlstCode"`
	STDSPCIESNM      string `gorm:"column:StdSpciesNm"`
	STDSPCIESCODE    string `gorm:"column:StdSpciesCode"`
	EXAMINUNITNM     string `gorm:"column:ExaminUnitNm"`
	EXAMINUNIT       string `gorm:"column:ExaminUnit"`
	STDUNITNM        string `gorm:"column:StdUnitNm"`
	STDUNITCODE      string `gorm:"column:StdUnitCode"`
	EXAMINGRADNM     string `gorm:"column:ExaminGradNm"`
	EXAMINGRADCODE   string `gorm:"column:ExaminGradCode"`
	STDGRADNM        string `gorm:"column:StdGradNm"`
	STDGRADCODE      string `gorm:"column:StdGradCode"`
	TODAYPRIC        string `gorm:"column:TodayPric"`
	BFRTPRIC         string `gorm:"column:BfrtPric"`
	IMPTRADE         string `gorm:"column:ImpTrade"`
	TRADEAMT         string `gorm:"column:TradeAmt"`
}

// MafraRetailPrice : mafra 소매시장 조사가격 테이블
type MafraRetailPrice struct {
	ExaminDate        string `gorm:"column:ExaminDate"`
	ExaminItemName    string `gorm:"column:ExaminItemName"`
	ExaminItemCode    string `gorm:"column:ExaminItemCode"`
	ExaminSpeciesName string `gorm:"column:ExaminSpeciesName"`
	ExaminSpeciesCode string `gorm:"column:ExaminSpeciesCode"`
	ExaminUnitName    string `gorm:"column:ExaminUnitName"`
	ExaminUnit        string `gorm:"column:ExaminUnit"`
	ExaminGradeName   string `gorm:"column:ExaminGradeName"`
	ExaminGradeCode   string `gorm:"column:ExaminGradeCode"`
	MinPrice          string `gorm:"column:MinPrice"`
	MaxPrice          string `gorm:"column:MaxPrice"`
}

// MafraWholePrice : mafra 도매시장 조사가격 테이블
type MafraWholePrice struct {
	ExaminDate        string `gorm:"column:ExaminDate"`
	ExaminItemName    string `gorm:"column:ExaminItemName"`
	ExaminItemCode    string `gorm:"column:ExaminItemCode"`
	ExaminSpeciesName string `gorm:"column:ExaminSpeciesName"`
	ExaminSpeciesCode string `gorm:"column:ExaminSpeciesCode"`
	ExaminUnitName    string `gorm:"column:ExaminUnitName"`
	ExaminUnit        string `gorm:"column:ExaminUnit"`
	ExaminGradeName   string `gorm:"column:ExaminGradeName"`
	ExaminGradeCode   string `gorm:"column:ExaminGradeCode"`
	Price             string `gorm:"column:Price"`
}

// TableName :mafra 농수축산 유통정보 조사가격(농수축산물표준코드변환) 테이블 명
func (MafraExamination) TableName() string {
	return "MAFRA_EXAMINATION"
}

// Columns :mafra 농수축산 유통정보 조사가격(농수축산물표준코드변환) 컬럼 명
func (MafraExamination) Columns() []string {
	return structs.Names(&MafraExamination{})
}

// TableName :mafra 소매시장 조사가격 테이블 명
func (MafraRetailPrice) TableName() string {
	return "MAFRA_RETAIL_PRICE"
}

// Columns :mafra 소매시장 조사가격 컬럼 명
func (MafraRetailPrice) Columns() []string {
	return structs.Names(&MafraRetailPrice{})
}

// TableName :mafra 도매시장 조사가격 테이블 명
func (MafraWholePrice) TableName() string {
	return "MAFRA_WHOLE_PRICE"
}

// Columns :mafra 도매시장 조사가격 컬럼 명
func (MafraWholePrice) Columns() []string {
	return structs.Names(&MafraWholePrice{})
}
