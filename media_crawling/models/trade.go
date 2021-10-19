package models

import (
	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

// Exportation : 정제 후 수출 데이터
type Exportation struct {
	HskPrdlstCode string `gorm:"column:HskPrdlstCode"`
	Weight        string `gorm:"column:Weight"`
	Amount        string `gorm:"column:Amount"`
	BaseDate      string `gorm:"column:BaseDate"`
}

// Importation : 정제 후 수입 데이터
type Importation struct {
	HskPrdlstCode string `gorm:"column:HskPrdlstCode"`
	Weight        string `gorm:"column:Weight"`
	Amount        string `gorm:"column:Amount"`
	BaseDate      string `gorm:"column:BaseDate"`
}

// Trade : 표준품목코드와 매핑되는 수출입 데이터
type Trade struct {
	HskPrdlstCode string `gorm:"column:HskPrdlstCode"`
	BaseDate      string `gorm:"column:BaseDate"`
	TradeType     string `gorm:"column:TradeType"`
	Weight        string `gorm:"column:Weight"`
	Amount        string `gorm:"column:Amount"`
}

// TableName : 수출 테이블 명
func (Exportation) TableName() string {
	return "AGRI_EXPORTATION"
}

// Columns : 수출 컬럼 명
func (Exportation) Columns() []string {
	return structs.Names(&Exportation{})
}

// GetExportation : Exportation 테이블 rows 획득
func GetExportation(db *gorm.DB) []Exportation {
	var exportations []Exportation
	db.Find(&exportations)

	return exportations
}

// TableName : 수입 테이블 명
func (Importation) TableName() string {
	return "AGRI_IMPORTATION"
}

// Columns : 수입 컬럼 명
func (Importation) Columns() []string {
	return structs.Names(&Importation{})
}

// GetImportation : Importation 테이블 rows 획득
func GetImportation(db *gorm.DB) []Importation {
	var importations []Importation
	db.Find(&importations)

	return importations
}

// TableName : 수입 테이블 명
func (Trade) TableName() string {
	return "TRADE"
}

// Columns : 수입 컬럼 명
func (Trade) Columns() []string {
	return structs.Names(&Trade{})
}
