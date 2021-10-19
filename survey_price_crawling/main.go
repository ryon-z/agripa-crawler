package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"crawler/survey_price_crawling/common"
	"crawler/survey_price_crawling/models"
)

var startTime time.Time
var issueList []string

func main() {
	startTime = time.Now()

	log.Printf("프로그램 시작")

	common.DBConnect("mysql", "agripa-collection-db.cn3ppiqdwmg7.ap-northeast-2.rds.amazonaws.com", "AGRIPA_COLLECTION", "agripa_dev", "agripa_d2v!")

	defer common.DBDisConnect()

	if len(os.Args) == 1 { //스케줄링
		log.Printf("[%s] 데이터를 수집합니다.\n", startTime.Format("2006-01-02"))

		scrapSurveyPrice(startTime.Format("20060102"))

	} else if len(os.Args) == 2 && len(os.Args[1]) == 10 { //특정일(2020-01-02)
		log.Printf("[%s] 데이터를 수집합니다.\n", os.Args[1])

		scrapSurveyPrice(os.Args[1])

	} else if len(os.Args) == 3 && len(os.Args[1]) == 10 && len(os.Args[2]) == 10 { //기간(2020-01-01 2020-12-31)
		log.Printf("[%s ~ %s] 데이터를 수집합니다.\n", os.Args[1], os.Args[2])

		stime, err := time.Parse("2006-01-02", os.Args[1])
		if err != nil {
			log.Panic(err)
		}

		etime, err := time.Parse("2006-01-02", os.Args[2])
		if err != nil {
			log.Panic(err)
		}

		for stime.Before(etime) || stime.Equal(etime) {

			log.Print(etime.Format("2006-01-02"))

			scrapSurveyPrice(etime.Format("20060102"))

			log.Printf("실행시간 : %s", time.Since(startTime))

			etime = etime.AddDate(0, 0, -1)
		}

	} else {
		log.Printf("허용되지 않는 인자 형식입니다.(%s)\n", os.Args)
	}

	log.Printf("실행시간 : %s", time.Since(startTime))
	log.Print("이슈날짜 트래킹 : ", issueList)
}

func scrapSurveyPrice(d string) {

	dd := strings.ReplaceAll(d, "-", "")

	dbCnt := models.GetSurveyPriceCount(dd)

	apiKey := "Id1ZE1pa3P7eFM%2FTv4mdg%2BE4fWstyv7CB7YOsOfROqlN6tibMq6r3mDJwYN1g2GulM2kb7F3jZB85L6IbuLViw%3D%3D"
	apikey2 := "XUa4iJfKVP2959o7OflfARznEV6jdarqNyT5gmW3A%2B77lEFf9U%2FQrwtEyBkCbKqulN7HaaCiLSrd5EiMMKdswA%3D%3D"

	apiQuery := fmt.Sprintf("http://apis.data.go.kr/B552895/LocalGovPriceInfoService/getSourcePriceResearchSearch?ServiceKey=%s&pageNo=%d&numOfRows=%d&_returnType=json&examin_de=%s",
		apiKey, 1, 10000, dd)

	apiQuery2 := fmt.Sprintf("http://apis.data.go.kr/B552895/LocalGovPriceInfoService/getSourcePriceResearchSearch?ServiceKey=%s&pageNo=%d&numOfRows=%d&_returnType=json&examin_de=%s",
		apikey2, 1, 10000, dd)

	log.Print(apiQuery)

	res, err := http.Get(apiQuery)
	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	var info models.LocalGovPriceInfo

	json.Unmarshal([]byte(data), &info)

	log.Printf("레코드 수 : %d", info.TotalCount)

	if info.TotalCount < 1 {

		// 임시 키로 다시 확인
		log.Print("임시키 호출")
		log.Print(apiQuery2)

		res2, err2 := http.Get(apiQuery2)
		if err2 != nil {
			log.Panic(err2)
		}

		defer res2.Body.Close()

		data2, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			log.Panic(err2)
		}

		json.Unmarshal([]byte(data2), &info)

		log.Printf("임시키 레코드 수 : %d", info.TotalCount)

		if info.TotalCount < 1 {
			issueList = append(issueList, d)

			return
		}
	}

	if dbCnt == info.TotalCount {
		log.Printf("==== DB SKIP (%s) ====", d)
		time.Sleep(time.Second * 1)

		return
	}

	for _, item := range info.List {
		models.AddSurveyPriceItem(item)
	}
}
