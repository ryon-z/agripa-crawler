package models

import (
	"crawler/survey_price_crawling/common"
)

// LocalGovPriceInfo 지자체 농수축산물 물가조사 가격정보
type LocalGovPriceInfo struct {
	TotalCount int                 `json:"totalCount"`
	List       []LocalGovPriceItem `json:"list"`
}

// LocalGovPriceItem 지자체 농수축산물 물가조사 가격정보 row
type LocalGovPriceItem struct {
	ExaminDe       string `json:"examin_de" gorm:"column:examin_de"`
	ExaminAreaNm   string `json:"examin_area_nm" gorm:"column:examin_area_nm"`
	ExaminAreaCd   string `json:"examin_area_cd" gorm:"column:examin_area_cd"`
	ExaminMrktNm   string `json:"examin_mrkt_nm" gorm:"column:examin_mrkt_nm"`
	ExaminMrktCd   string `json:"examin_mrkt_cd" gorm:"column:examin_mrkt_cd"`
	PrdlstDetailNm string `json:"prdlst_detail_nm" gorm:"column:prdlst_detail_nm"`
	PrdlstCd       string `json:"prdlst_cd" gorm:"column:prdlst_cd"`
	DistbStepSe    string `json:"distb_step_se" gorm:"column:distb_step_se"`
	DistbStep      string `json:"distb_step" gorm:"column:distb_step"`
	PrdlstDetailCd string `json:"prdlst_detail_cd" gorm:"column:prdlst_detail_cd"`
	GradCd         string `json:"grad_cd" gorm:"column:grad_cd"`
	PrdlstNm       string `json:"prdlst_nm" gorm:"column:prdlst_nm"`
	ExaminAmt      string `json:"examin_amt" gorm:"column:examin_amt"`
	BfrtExaminAmt  string `json:"bfrt_examin_amt" gorm:"column:bfrt_examin_amt"`
	Stndrd         string `json:"stndrd" gorm:"column:stndrd"`
	Grad           string `json:"grad" gorm:"column:grad"`
}

// TableName 조사가격정보 테이블 명
func (LocalGovPriceItem) TableName() string {
	return "SURVEY_PRICE"
}

// AddSurveyPriceItem 조사가격 정보 입력
func AddSurveyPriceItem(row LocalGovPriceItem) int {

	db := common.GetDB()

	db.Create(row)

	return 0
}

// GetSurveyPriceCount 조사가격 카운트 조회
func GetSurveyPriceCount(d string) int {
	db := common.GetDB()

	var cnt int

	db.Table("SURVEY_PRICE").Where("examin_de = ?", d).Count(&cnt)

	return cnt
}
