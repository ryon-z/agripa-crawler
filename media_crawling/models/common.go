package models

import (
	"reflect"
)

// Model : 범용 모델 interface
type Model interface {
	TableName() string
	Columns() []string
}

var newsArr []News
var blogArr []Blog
var pressArr []Press
var queryArr []Query
var cafeArr []Cafe
var youtubeArr []Youtube
var channelArr []Channel
var realtimeAuctionArr []RealtimeAuction
var adjustedAuctionArr []AdjustedAuction
var mapStdExamItemArr []MapStdExamItem
var importationArr []Importation
var exportationArr []Exportation
var weatherArr []Weather
var garakCodeArr []GarakCode
var wholesaleMarketCodeArr []WholesaleMarketCode
var wholesaleMarketCoCodeArr []WholesaleMarketCoCode
var stdGradeCodeArr []StdGradeCode
var stdUnitCodeArr []StdUnitCode
var placeOriginCodeArr []PlaceOriginCode
var rawBreakingAuctionArr []RawBreakingAuction
var stdSpeciesCodeArr []StdSpeciesCode
var breakingAuctionArr []BreakingAuction
var mafraAdjAuctionStatsArr []MafraAdjAuctionStats
var mafraExaminationArr []MafraExamination
var mafraRetailPriceArr []MafraRetailPrice
var mafraWholePriceArr []MafraWholePrice
var tradeArr []Trade
var mafraAdjAuctionQuantityArr []MafraAdjAuctionQuantity
var stdItemKeywordArr []StdItemKeyword

var modelArrAddrMap map[string]interface{}

// GetModel : DB 구조체를 리턴한다.
// (허용 값 : news, blog, press, query)
func GetModel(modelName string) Model {
	modelMap := map[string]Model{
		"news":                    News{},
		"blog":                    Blog{},
		"press":                   Press{},
		"query":                   Query{},
		"cafe":                    Cafe{},
		"youtube":                 Youtube{},
		"channel":                 Channel{},
		"realtimeAuction":         RealtimeAuction{},
		"adjustedAuction":         AdjustedAuction{},
		"mapStdExamItem":          MapStdExamItem{},
		"importation":             Importation{},
		"exportation":             Exportation{},
		"weather":                 Weather{},
		"garakCode":               GarakCode{},
		"wholesaleMarketCode":     WholesaleMarketCode{},
		"wholesaleMarketCoCode":   WholesaleMarketCoCode{},
		"stdGradeCode":            StdGradeCode{},
		"stdUnitCode":             StdUnitCode{},
		"placeOriginCode":         PlaceOriginCode{},
		"rawBreakingAuction":      RawBreakingAuction{},
		"stdSpeciesCode":          StdSpeciesCode{},
		"breakingAuction":         BreakingAuction{},
		"mafraAdjAuctionStats":    MafraAdjAuctionStats{},
		"mafraExamination":        MafraExamination{},
		"mafraRetailPrice":        MafraRetailPrice{},
		"mafraWholePrice":         MafraWholePrice{},
		"trade":                   Trade{},
		"mafraAdjAuctionQuantity": MafraAdjAuctionQuantity{},
		"stdItemKeyword":          StdItemKeyword{},
	}

	return modelMap[modelName]
}

// GetModelArrAddr : 모델 배열 주소 값 획득
func GetModelArrAddr(modelName string) interface{} {
	modelArrAddrMap := map[string]interface{}{
		"news":                    &newsArr,
		"blog":                    &blogArr,
		"press":                   &pressArr,
		"query":                   &queryArr,
		"cafe":                    &cafeArr,
		"youtube":                 &youtubeArr,
		"channel":                 &channelArr,
		"realtimeAuction":         &realtimeAuctionArr,
		"adjustedAuction":         &adjustedAuctionArr,
		"mapStdExamItem":          &mapStdExamItemArr,
		"importation":             &importationArr,
		"exportation":             &exportationArr,
		"weather":                 &weatherArr,
		"garakCode":               &garakCodeArr,
		"wholesaleMarketCode":     &wholesaleMarketCodeArr,
		"wholesaleMarketCoCode":   &wholesaleMarketCoCodeArr,
		"stdGradeCode":            &stdGradeCodeArr,
		"stdUnitCode":             &stdUnitCodeArr,
		"placeOriginCode":         &placeOriginCodeArr,
		"rawBreakingAuction":      &rawBreakingAuctionArr,
		"stdSpeciesCode":          &stdSpeciesCodeArr,
		"breakingAuction":         &breakingAuctionArr,
		"mafraRetailPrice":        &mafraRetailPriceArr,
		"mafraWholePrice":         &mafraWholePriceArr,
		"trade":                   &tradeArr,
		"mafraAdjAuctionQuantity": &mafraAdjAuctionQuantityArr,
		"stdItemKeyword":          &stdItemKeywordArr,
	}

	return modelArrAddrMap[modelName]
}

// GetInterfaceArr : 입력 받은 모델 배열 주소 안에 값을 []interface{} 형으로 변환
func GetInterfaceArr(modelArrAddr interface{}) []interface{} {
	val := reflect.ValueOf(modelArrAddr).Elem()
	data := make([]interface{}, val.Len())
	for i := 0; i < val.Len(); i++ {
		data[i] = val.Index(i).Interface()
	}

	return data
}
