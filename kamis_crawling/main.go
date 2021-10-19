package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	m "kamis_crawling/Models"
	"kamis_crawling/common"
)

var startTime time.Time

func main() {
	startTime = time.Now()

	if len(os.Args) < 2 {
		panic("명령 인수가 부족합니다.")
	}

	common.DBConnect("mysql", "agripa-operation-db.cn3ppiqdwmg7.ap-northeast-2.rds.amazonaws.com", "AGRIPA", "agripa_dev", "agripa_d2v!")

	defer common.DBDisConnect()

	var strClass string

	class := os.Args[1]

	if class == "-whole" {
		strClass = "02"
	} else if class == "-retail" {
		strClass = "01"
	} else {
		panic("첫번째 인자 값이 잘 못 되었습니다.(-whole, -retail)")
	}

	if len(os.Args) == 2 { // daily
		scrapPeriodPrice(strClass, startTime.Format("2006-01-02"), startTime.Format("2006-01-02"))

		fmt.Println("총 실행 시간: ", time.Since(startTime))

		return
	}

	if len(os.Args) == 3 && len(os.Args[2]) == 10 { //특정 날짜
		scrapPeriodPrice(strClass, os.Args[2], os.Args[2])

		fmt.Println("총 실행 시간: ", time.Since(startTime))

		return
	}

	if len(os.Args[2]) != 7 || len(os.Args[3]) != 7 { // 월 단위 범위
		panic("기간 인자가 잘 못 되었습니다.(2006-01)")
	}

	stime, err := time.Parse("2006-01-02", os.Args[2]+"-01")
	if err != nil {
		panic(err)
	}

	etime, err := time.Parse("2006-01-02", os.Args[3]+"-01")
	if err != nil {
		panic(err)
	}

	for stime.Before(etime) || stime.Equal(etime) {

		scrapPeriodPrice(strClass, etime.Format("2006-01-02"), etime.AddDate(0, 1, -1).Format("2006-01-02"))

		etime = etime.AddDate(0, -1, 0)
	}

	fmt.Println("총 실행 시간: ", time.Since(startTime))

	return
}

func scrapPeriodPrice(cls string, sdate string, edate string) {

	codelist := m.GetItemCodeList(cls)

	for _, pcode := range codelist {
		fmt.Printf("%+v\t%s,%s\n", pcode, sdate, edate)

		apiQuery := fmt.Sprintf("http://www.kamis.or.kr/service/price/xml.do?action=periodProductList&p_productclscode=%s&p_startday=%s&p_endday=%s&p_itemcode=%d&p_kindcode=%s&p_productrankcode=%s&p_cert_key=6111613e-b52f-47c8-87ec-c405d564506c&p_cert_id=dev@pandac.co.kr&p_returntype=json",
			cls, sdate, edate, pcode.ItemCode, pcode.ItemKindCode, pcode.GradeCode)

		// GET 호출
		resp, err := http.Get(apiQuery)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		// 결과 출력
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var kdata m.KamisApi
		json.Unmarshal([]byte(data), &kdata)

		fmt.Print("조회 레코드 수: ", len(kdata.Data.Item))

		for _, item := range kdata.Data.Item {
			var amt float64
			var unit string

			if item.Itemname != "" {

				arrstr := regexp.MustCompile("[()]").Split(item.Kindname, 10)

				for _, str := range arrstr {
					if regexp.MustCompile("[0-9.]{1,}").MatchString(str) {
						strAmt := regexp.MustCompile("[0-9.]{1,}").FindString(str)
						amt, err = strconv.ParseFloat(strAmt, 32)
						if err != nil {
							//fmt.Println(err)
							amt = 0
						}

						unit = strings.ReplaceAll(str, strAmt, "")
					}
				}

				price, err := strconv.Atoi(strings.ReplaceAll(item.Price, ",", ""))
				if err != nil {
					//fmt.Println(err)
					price = 0
				}

				if cls == "01" {
					var newRow m.RetailPrice
					newRow.ShipDate = fmt.Sprintf("%s-%s", item.Yyyy, strings.ReplaceAll(item.Regday, "/", "-"))
					newRow.ItemCode = pcode.ItemCode
					newRow.ItemKindCode = pcode.ItemKindCode
					newRow.GradeCode = pcode.GradeCode
					newRow.MarketName = item.Marketname
					newRow.AreaName = item.Countyname
					newRow.ShipPrice = price
					newRow.ShipUnit = unit
					newRow.ShipAmt = amt

					//db.Create(newRow)
					m.AddRetailPrice(newRow)
				} else if cls == "02" {
					var newRow m.WholePrice
					newRow.ShipDate = fmt.Sprintf("%s-%s", item.Yyyy, strings.ReplaceAll(item.Regday, "/", "-"))
					newRow.ItemCode = pcode.ItemCode
					newRow.ItemKindCode = pcode.ItemKindCode
					newRow.GradeCode = pcode.GradeCode
					newRow.MarketName = item.Marketname
					newRow.AreaName = item.Countyname
					newRow.ShipPrice = price
					newRow.ShipUnit = unit
					newRow.ShipAmt = amt

					//db.Create(newRow)
					m.AddWholePrice(newRow)
				}
			}
		}

		fmt.Println("\t실행시간: ", time.Since(startTime))

		time.Sleep(time.Second * 1)
	}
}
