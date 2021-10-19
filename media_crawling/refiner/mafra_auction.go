package refiner

import (
	"errors"
	"fmt"
	"media_crawling/alarm"
	"media_crawling/crawler"
	"media_crawling/models"
	"media_crawling/util"
	"time"
)

// RefineMafraAcution : MAFRA 경락가격 정제 후 DB import
func RefineMafraAcution(modelName string, start time.Time, end time.Time) {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	if !util.InArray(modelName, []string{"mafraAdjAuctionStats"}) {
		errorMessage := fmt.Sprintf("REFINE_MAFRA_AUCTION :: 잘못된 modelName, modelName: %s", modelName)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}

	model := models.GetModel(modelName)
	rawTableName := model.TableName()

	wherePhrases := GetDateWherePhrase("AucngDe", "", start, end)
	for i, wherePhrase := range wherePhrases {
		year := start.Year() + i
		tableName := fmt.Sprintf("%s_%d", rawTableName, year)
		sqlQuery := fmt.Sprintf(`
			INSERT IGNORE INTO %s
			SELECT A.* 
			FROM %s AS A
			INNER JOIN (SELECT DISTINCT StdSpciesCode from MAP_STD_EXAM_ITEM) AS B
			ON A.SpciesCd = B.StdSpciesCode
			%s
		;`, rawTableName, tableName, wherePhrase)
		crawler.RunQuery(db, sqlQuery)
	}
}

func UpdateMafraAdjAuctionQuantity(start time.Time, end time.Time) {
	db := crawler.GetDB("collection")
	defer crawler.CloseDB(db)

	auctionQuantityTableName := models.MafraAdjAuctionQuantity{}.TableName()
	unitConvertingTableName := "STD_UNIT_CONVERTING"

	crawler.IsTableExists(db, "AGRIPA_COLLECTION", models.MafraAdjAuctionStats{}.TableName())
	crawler.IsTableExists(db, "AGRIPA_COLLECTION", unitConvertingTableName)
	crawler.IsTableExists(db, "AGRIPA_COLLECTION", auctionQuantityTableName)

	auctionStatsTableName := models.MafraAdjAuctionStats{}.TableName()
	startDate := util.GetDateString(start, "")
	endDate := util.GetDateString(end, "")

	sqlQuery := fmt.Sprintf(`
		INSERT IGNORE INTO %s
		SELECT date_format(str_to_date(AucngDe, '%%Y%%m%%d'),'%%Y-%%m-%%d'), PrdlstNm, PrdlstCd, SpciesNm, SpciesCd, TRUNCATE(acc_qy * ConvertedValue, 3)
		FROM (
			SELECT *
			FROM (
				SELECT AucngDe, PrdlstNm, PrdlstCd, SpciesNm, SpciesCd, StndrdCd, acc_qy,
						@qy_rank := IF(@curr_AucngDe = AucngDe AND @curr_PrdlstCd = PrdlstCd, @qy_rank+1, 1) AS qy_rank,
						@curr_AucngDe := AucngDe,
						@curr_PrdlstCd := PrdlstCd
				FROM (
					SELECT AucngDe, PrdlstNm, PrdlstCd, SpciesNm, SpciesCd, StndrdCd, sum(DelngQy) as acc_qy
					FROM %s, (SELECT @qy_rank, @curr_AucngDe, @curr_PrdlstCd) X
					WHERE AucngDe >= %s
					AND AucngDe <= %s 
					GROUP BY AucngDe, SpciesCd
					ORDER BY AucngDe ASC, PrdlstCd ASC, acc_qy DESC
				) AS TMP
			) AS RANKED
			WHERE qy_rank <=3
		) as A
		JOIN %s AS B
		ON A.StndrdCd = B.UnitCode
	;`, auctionQuantityTableName, auctionStatsTableName, startDate, endDate, unitConvertingTableName)
	crawler.RunQuery(db, sqlQuery)
}
