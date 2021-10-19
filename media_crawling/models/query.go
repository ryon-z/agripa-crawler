package models

import (
	"fmt"
	"media_crawling/util"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
)

// Query : 언론사 뉴스 검색어 구조체
type Query struct {
	ItemCode string `gorm:"column:ItemCode"`
	Query    string `gorm:"column:Query"`
	Usage    string `gorm:"column:Usage"`
	Priority string `gorm:"column:Priority;default:0"`
}

// TableName : 언론사 뉴스 검색어 테이블 명
func (Query) TableName() string {
	return "AGRI_QUERY"
}

// Columns : 언론사 뉴스 검색어 명
func (Query) Columns() []string {
	return structs.Names(&Query{})
}

// GetQueries : 네이버 API 수집 시 검색할 검색어 획득
func GetQueries(db *gorm.DB, usage string, firstCodeNum int) []string {
	var queries []Query
	var selectedQueries map[string]Query
	selectedQueries = make(map[string]Query)
	var results []string
	itemCodeWherePhrase := ""

	if firstCodeNum > 0 && firstCodeNum <= 6 {
		itemCodeWherePhrase = fmt.Sprintf("AND itemCode like '%d%%'", firstCodeNum)
	}

	tableName := Query{}.TableName()
	sqlQuery := fmt.Sprintf("SELECT * FROM %s WHERE `usage` = '%s' %s", tableName, usage, itemCodeWherePhrase)
	db.Raw(sqlQuery).Scan(&queries)

	// 품목코드(ItemCode) 별 우선순위 높은 queries만 추출
	for _, query := range queries {
		key := query.ItemCode
		if val, ok := selectedQueries[key]; ok {
			if query.Priority > val.Priority {
				selectedQueries[key] = query
			}
		} else {
			selectedQueries[key] = query
		}
	}
	fmt.Println("품목코드(ItemCode) 별 우선순위 높은 newsQuries만 추출")
	fmt.Println(selectedQueries)
	fmt.Println("총 개수", len(selectedQueries))
	fmt.Println()

	// 배열로 변환 후 리턴
	for _, query := range selectedQueries {
		if !util.InArray(query.Query, results) {
			results = append(results, query.Query)
		}
	}
	fmt.Println("유니크한 queries만 모아 배열로 변환")
	fmt.Println(results)
	fmt.Println("총 개수", len(results))
	fmt.Println()

	return results
}
