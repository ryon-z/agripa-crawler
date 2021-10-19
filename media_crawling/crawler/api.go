package crawler

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"media_crawling/alarm"
	"media_crawling/config"
	"media_crawling/util"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// Request : api 요청 시 사용하는 범용 Request 함수
func Request(url string, headers map[string]string) string {
	// Request 객체 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}

	//필요시 헤더 추가 가능
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Client객체에서 Request 실행
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		alarm.PostMessage("default", err.Error())
		panic(err)
	}
	defer resp.Body.Close()

	// 결과 출력
	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes) //바이트를 문자열로

	return str
}

// GetUselessWords : 불필요한 단어와 해당 단어의 대체값을 가진 map을 리턴
func GetUselessWords() map[string]string {
	uselessWords := map[string]string{
		"<b>":    "",
		"</b>":   "",
		"&quot;": "",
		"&amp;":  "",
		"\"":     "'",
		"\\":     "",
		"&gt;":   "",
		"&lt;":   "",
		"&#39;":  "'",
	}

	return uselessWords
}

// NaverNewsResponse : 네이버 뉴스 응답값 구조체
type NaverNewsResponse struct {
	LastBuildDate string `json:"lastBuildDate"`
	Total         int    `json:"total"`
	Start         int    `json:"start"`
	Display       int    `json:"display"`
	Items         []struct {
		Title        string `json:"title"`
		Originallink string `json:"originallink"`
		Link         string `json:"link"`
		Description  string `json:"description"`
		PubDate      string `json:"pubDate"`
	} `json:"items"`
}

// NaverBlogResponse : 네이버 블로그 응답값 구조체
type NaverBlogResponse struct {
	LastBuildDate string `json:"lastBuildDate"`
	Total         int    `json:"total"`
	Start         int    `json:"start"`
	Display       int    `json:"display"`
	Items         []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		Description string `json:"description"`
		Bloggername string `json:"bloggername"`
		Bloggerlink string `json:"bloggerlink"`
		Postdate    string `json:"postdate"`
	} `json:"items"`
}

// NaverCafeResponse : 네이버 카페 응답값 구조체
type NaverCafeResponse struct {
	LastBuildDate string `json:"lastBuildDate"`
	Total         int    `json:"total"`
	Start         int    `json:"start"`
	Display       int    `json:"display"`
	Items         []struct {
		Title       string `json:"title"`
		Link        string `json:"link"`
		Description string `json:"description"`
		Cafename    string `json:"cafename"`
		Cafeurl     string `json:"cafeurl"`
	} `json:"items"`
}

// YoutubeResponse : 유튜브 검색 응답값 구조체
type YoutubeResponse struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt string `json:"publishedAt"`
			ChannelID   string `json:"channelId"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
			ChannelTitle         string `json:"channelTitle"`
			LiveBroadcastContent string `json:"liveBroadcastContent"`
			PublishTime          string `json:"publishTime"`
		} `json:"snippet"`
	} `json:"items"`
}

// AuctionResponse : 경락가격 응답값 구조체
type AuctionResponse struct {
	List []struct {
		ReturnType       string `json:"_returnType"`
		CprCode          string `json:"cpr_code"`
		CprInsttCode     string `json:"cpr_instt_code"`
		CprInsttNm       string `json:"cpr_instt_nm"`
		CprNm            string `json:"cpr_nm"`
		DelngDe          string `json:"delng_de"`
		DelngPric        string `json:"delng_pric"`
		DelngPrutCode    string `json:"delng_prut_code"`
		DelngQy          string `json:"delng_qy"`
		LclasCode        string `json:"lclas_code"`
		LclasNm          string `json:"lclas_nm"`
		MlsfcCode        string `json:"mlsfc_code"`
		MlsfcNm          string `json:"mlsfc_nm"`
		NumOfRows        string `json:"numOfRows"`
		PageNo           string `json:"pageNo"`
		ResultCode       string `json:"resultCode"`
		ResultMsg        string `json:"resultMsg"`
		SbidPricAvg      string `json:"sbid_pric_avg"`
		SbidPricMax      string `json:"sbid_pric_max"`
		SbidPricMin      string `json:"sbid_pric_min"`
		SbidPricMvAvg    string `json:"sbid_pric_mv_avg"`
		ServiceKey       string `json:"serviceKey"`
		StdFrmlcNewCode  string `json:"std_frmlc_new_code"`
		StdFrmlcNewNm    string `json:"std_frmlc_new_nm"`
		StdGradCode      string `json:"std_grad_code"`
		StdGradNm        string `json:"std_grad_nm"`
		StdMgNewCode     string `json:"std_mg_new_code"`
		StdMgNewNm       string `json:"std_mg_new_nm"`
		StdMtcCode       string `json:"std_mtc_code"`
		StdMtcNewCode    string `json:"std_mtc_new_code"`
		StdMtcNewNm      string `json:"std_mtc_new_nm"`
		StdMtcNm         string `json:"std_mtc_nm"`
		StdPrdlstCode    string `json:"std_prdlst_code"`
		StdPrdlstNewCode string `json:"std_prdlst_new_code"`
		StdPrdlstNewNm   string `json:"std_prdlst_new_nm"`
		StdPrdlstNm      string `json:"std_prdlst_nm"`
		StdQlityNewCode  string `json:"std_qlity_new_code"`
		StdQlityNewNm    string `json:"std_qlity_new_nm"`
		StdUnitCode      string `json:"std_unit_code"`
		StdUnitNewCode   string `json:"std_unit_new_code"`
		StdUnitNewNm     string `json:"std_unit_new_nm"`
		StdUnitNm        string `json:"std_unit_nm"`
		TotalCount       string `json:"totalCount"`
		WhsalCode        string `json:"whsal_code"`
		WhsalMrktCode    string `json:"whsal_mrkt_code"`
		WhsalMrktNm      string `json:"whsal_mrkt_nm"`
		WhsalNm          string `json:"whsal_nm"`
	} `json:"list"`
	Parm struct {
		ReturnType       string `json:"_returnType"`
		CprCode          string `json:"cpr_code"`
		CprInsttCode     string `json:"cpr_instt_code"`
		CprInsttNm       string `json:"cpr_instt_nm"`
		CprNm            string `json:"cpr_nm"`
		DelngDe          string `json:"delng_de"`
		DelngPric        string `json:"delng_pric"`
		DelngPrutCode    string `json:"delng_prut_code"`
		DelngQy          string `json:"delng_qy"`
		LclasCode        string `json:"lclas_code"`
		LclasNm          string `json:"lclas_nm"`
		MlsfcCode        string `json:"mlsfc_code"`
		MlsfcNm          string `json:"mlsfc_nm"`
		NumOfRows        string `json:"numOfRows"`
		PageNo           string `json:"pageNo"`
		ResultCode       string `json:"resultCode"`
		ResultMsg        string `json:"resultMsg"`
		SbidPricAvg      string `json:"sbid_pric_avg"`
		SbidPricMax      string `json:"sbid_pric_max"`
		SbidPricMin      string `json:"sbid_pric_min"`
		SbidPricMvAvg    string `json:"sbid_pric_mv_avg"`
		ServiceKey       string `json:"serviceKey"`
		StdFrmlcNewCode  string `json:"std_frmlc_new_code"`
		StdFrmlcNewNm    string `json:"std_frmlc_new_nm"`
		StdGradCode      string `json:"std_grad_code"`
		StdGradNm        string `json:"std_grad_nm"`
		StdMgNewCode     string `json:"std_mg_new_code"`
		StdMgNewNm       string `json:"std_mg_new_nm"`
		StdMtcCode       string `json:"std_mtc_code"`
		StdMtcNewCode    string `json:"std_mtc_new_code"`
		StdMtcNewNm      string `json:"std_mtc_new_nm"`
		StdMtcNm         string `json:"std_mtc_nm"`
		StdPrdlstCode    string `json:"std_prdlst_code"`
		StdPrdlstNewCode string `json:"std_prdlst_new_code"`
		StdPrdlstNewNm   string `json:"std_prdlst_new_nm"`
		StdPrdlstNm      string `json:"std_prdlst_nm"`
		StdQlityNewCode  string `json:"std_qlity_new_code"`
		StdQlityNewNm    string `json:"std_qlity_new_nm"`
		StdUnitCode      string `json:"std_unit_code"`
		StdUnitNewCode   string `json:"std_unit_new_code"`
		StdUnitNewNm     string `json:"std_unit_new_nm"`
		StdUnitNm        string `json:"std_unit_nm"`
		TotalCount       string `json:"totalCount"`
		WhsalCode        string `json:"whsal_code"`
		WhsalMrktCode    string `json:"whsal_mrkt_code"`
		WhsalMrktNm      string `json:"whsal_mrkt_nm"`
		WhsalNm          string `json:"whsal_nm"`
	} `json:"parm"`
	TotalCount          int `json:"totalCount"`
	RealtimeMktStatsSVO struct {
		ReturnType       string `json:"_returnType"`
		CprCode          string `json:"cpr_code"`
		CprInsttCode     string `json:"cpr_instt_code"`
		CprInsttNm       string `json:"cpr_instt_nm"`
		CprNm            string `json:"cpr_nm"`
		DelngDe          string `json:"delng_de"`
		DelngPric        string `json:"delng_pric"`
		DelngPrutCode    string `json:"delng_prut_code"`
		DelngQy          string `json:"delng_qy"`
		LclasCode        string `json:"lclas_code"`
		LclasNm          string `json:"lclas_nm"`
		MlsfcCode        string `json:"mlsfc_code"`
		MlsfcNm          string `json:"mlsfc_nm"`
		NumOfRows        string `json:"numOfRows"`
		PageNo           string `json:"pageNo"`
		ResultCode       string `json:"resultCode"`
		ResultMsg        string `json:"resultMsg"`
		SbidPricAvg      string `json:"sbid_pric_avg"`
		SbidPricMax      string `json:"sbid_pric_max"`
		SbidPricMin      string `json:"sbid_pric_min"`
		SbidPricMvAvg    string `json:"sbid_pric_mv_avg"`
		ServiceKey       string `json:"serviceKey"`
		StdFrmlcNewCode  string `json:"std_frmlc_new_code"`
		StdFrmlcNewNm    string `json:"std_frmlc_new_nm"`
		StdGradCode      string `json:"std_grad_code"`
		StdGradNm        string `json:"std_grad_nm"`
		StdMgNewCode     string `json:"std_mg_new_code"`
		StdMgNewNm       string `json:"std_mg_new_nm"`
		StdMtcCode       string `json:"std_mtc_code"`
		StdMtcNewCode    string `json:"std_mtc_new_code"`
		StdMtcNewNm      string `json:"std_mtc_new_nm"`
		StdMtcNm         string `json:"std_mtc_nm"`
		StdPrdlstCode    string `json:"std_prdlst_code"`
		StdPrdlstNewCode string `json:"std_prdlst_new_code"`
		StdPrdlstNewNm   string `json:"std_prdlst_new_nm"`
		StdPrdlstNm      string `json:"std_prdlst_nm"`
		StdQlityNewCode  string `json:"std_qlity_new_code"`
		StdQlityNewNm    string `json:"std_qlity_new_nm"`
		StdUnitCode      string `json:"std_unit_code"`
		StdUnitNewCode   string `json:"std_unit_new_code"`
		StdUnitNewNm     string `json:"std_unit_new_nm"`
		StdUnitNm        string `json:"std_unit_nm"`
		TotalCount       string `json:"totalCount"`
		WhsalCode        string `json:"whsal_code"`
		WhsalMrktCode    string `json:"whsal_mrkt_code"`
		WhsalMrktNm      string `json:"whsal_mrkt_nm"`
		WhsalNm          string `json:"whsal_nm"`
	} `json:"Realtime_mkt_stats_sVO"`
}

// WeatherResponse : 날씨 구조체 응답값 구조체
type WeatherResponse struct {
	Response struct {
		Header struct {
			ResultCode string `json:"resultCode"`
			ResultMsg  string `json:"resultMsg"`
		} `json:"header"`
		Body struct {
			DataType string `json:"dataType"`
			Items    struct {
				Item []struct {
					AreaID        string `json:"areaId"`
					AreaName      string `json:"areaName"`
					DayAvgRhm     int    `json:"dayAvgRhm"`
					DayAvgTa      int    `json:"dayAvgTa"`
					DayAvgWs      int    `json:"dayAvgWs"`
					DayMaxTa      int    `json:"dayMaxTa"`
					DayMinRhm     int    `json:"dayMinRhm"`
					DayMinTa      int    `json:"dayMinTa"`
					DaySumRn      int    `json:"daySumRn"`
					DaySumSs      int    `json:"daySumSs"`
					PaCropName    string `json:"paCropName"`
					PaCropSpeID   string `json:"paCropSpeId"`
					PaCropSpeName string `json:"paCropSpeName"`
					WrnCd         string `json:"wrnCd"`
					WrnCount      int    `json:"wrnCount"`
					Ymd           string `json:"ymd"`
				} `json:"item"`
			} `json:"items"`
			PageNo     int `json:"pageNo"`
			NumOfRows  int `json:"numOfRows"`
			TotalCount int `json:"totalCount"`
		} `json:"body"`
	} `json:"response"`
}

// GarakCodeResponse : 가락시장 코드 응답값 구조체
type GarakCodeResponse struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string `xml:",chardata"`
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Items struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text       string `xml:",chardata"`
				GarrakCode string `xml:"garrakCode"`
				GarrakName string `xml:"garrakName"`
				Rnum       string `xml:"rnum"`
				Sclasscode string `xml:"sclasscode"`
				StanCode   string `xml:"stanCode"`
			} `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

// WholesaleMarketCodeResponse : 도매시장 코드 응답값 구조체
type WholesaleMarketCodeResponse struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string `xml:",chardata"`
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Items struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text     string `xml:",chardata"`
				Marketco string `xml:"marketco"`
				Marketnm string `xml:"marketnm"`
				Rnum     string `xml:"rnum"`
			} `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

// WholesaleMarketCoCodeResponse : 도매시장 법인 코드 응답값 구조체
type WholesaleMarketCoCodeResponse struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string `xml:",chardata"`
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Items struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text       string `xml:",chardata"`
				Cocode     string `xml:"cocode"`
				Coname     string `xml:"coname"`
				Marketcode string `xml:"marketcode"`
				Marketname string `xml:"marketname"`
				Rnum       string `xml:"rnum"`
			} `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

// StdGradeCodeResponse : 표준 등급 코드 응답값 구조체
type StdGradeCodeResponse struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string `xml:",chardata"`
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Items struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text      string `xml:",chardata"`
				Gradecode string `xml:"gradecode"`
				Gradename string `xml:"gradename"`
				Rnum      string `xml:"rnum"`
			} `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

// StdUnitCodeResponse : 표준 단위 코드 응답값 구조체
type StdUnitCodeResponse struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string `xml:",chardata"`
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Items struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text     string `xml:",chardata"`
				Rnum     string `xml:"rnum"`
				Unitcode string `xml:"unitcode"`
				Unitname string `xml:"unitname"`
			} `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

// PlaceOriginCodeResponse : 산지 코드 응답값 구조체
type PlaceOriginCodeResponse struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string `xml:",chardata"`
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Items struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text    string `xml:",chardata"`
				Rnum    string `xml:"rnum"`
				Sido    string `xml:"sido"`
				Sigun   string `xml:"sigun"`
				Zipcode string `xml:"zipcode"`
				Dong    string `xml:"dong"`
			} `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

// BreakingAuctionResponse : 실시간 경락 속보 응답값 구조체
type BreakingAuctionResponse struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text       string `xml:",chardata"`
		ResultCode string `xml:"resultCode"`
		ResultMsg  string `xml:"resultMsg"`
	} `xml:"header"`
	Body struct {
		Text  string `xml:",chardata"`
		Items struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text       string `xml:",chardata"`
				Bidtime    string `xml:"bidtime"`
				Chulagtnm  string `xml:"chulagtnm"`
				Coname     string `xml:"coname"`
				Gradename  string `xml:"gradename"`
				Marketname string `xml:"marketname"`
				Mclassname string `xml:"mclassname"`
				Price      string `xml:"price"`
				Sanji      string `xml:"sanji"`
				Sclassname string `xml:"sclassname"`
				Tradeamt   string `xml:"tradeamt"`
				Unitname   string `xml:"unitname"`
			} `xml:"item"`
		} `xml:"items"`
		NumOfRows  string `xml:"numOfRows"`
		PageNo     string `xml:"pageNo"`
		TotalCount string `xml:"totalCount"`
	} `xml:"body"`
}

// StdSpeciesResponse : 표준품종코드 응답값 구조체
type StdSpeciesResponse struct {
	Grid201412210000000001201 struct {
		TotalCnt int `json:"totalCnt"`
		StartRow int `json:"startRow"`
		EndRow   int `json:"endRow"`
		Result   struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"result"`
		Row []struct {
			ROWNUM     int    `json:"ROW_NUM"`
			LClassCode string `json:"CATGORY_CD"`
			LClassName string `json:"CATGORY_NM"`
			MClassCode string `json:"PRDLST_CD"`
			MClassName string `json:"PRDLST_NM"`
			SClassCode string `json:"SPCIES_CD"`
			SClassName string `json:"SPCIES_NM"`
		} `json:"row"`
	} `json:"Grid_20141221000000000120_1"`
}

// MafraAdjAuctionStatsResponse : mafra 일별 정산 경락 요약정보 응답값 구조체
type MafraAdjAuctionStatsResponse struct {
	Grid201606240000000003481 struct {
		TotalCnt int `json:"totalCnt"`
		StartRow int `json:"startRow"`
		EndRow   int `json:"endRow"`
		Result   struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"result"`
		Row []struct {
			ROWNUM            int    `json:"ROW_NUM"`
			AUCNGDE           string `json:"AUCNG_DE"`
			PBLMNGWHSALMRKTNM string `json:"PBLMNG_WHSAL_MRKT_NM"`
			PBLMNGWHSALMRKTCD string `json:"PBLMNG_WHSAL_MRKT_CD"`
			CPRNM             string `json:"CPR_NM"`
			CPRCD             string `json:"CPR_CD"`
			PRDLSTNM          string `json:"PRDLST_NM"`
			PRDLSTCD          string `json:"PRDLST_CD"`
			SPCIESNM          string `json:"SPCIES_NM"`
			SPCIESCD          string `json:"SPCIES_CD"`
			DELNGBUNDLEQY     int    `json:"DELNGBUNDLE_QY"`
			STNDRD            string `json:"STNDRD"`
			STNDRDCD          string `json:"STNDRD_CD"`
			GRAD              string `json:"GRAD"`
			GRADCD            string `json:"GRAD_CD"`
			SANJICD           string `json:"SANJI_CD"`
			SANJINM           string `json:"SANJI_NM"`
			MUMMAMT           int    `json:"MUMM_AMT"`
			AVRGAMT           int    `json:"AVRG_AMT"`
			MXMMAMT           int    `json:"MXMM_AMT"`
			DELNGQY           int    `json:"DELNG_QY"`
			CNTS              int    `json:"CNTS"`
		} `json:"row"`
	} `json:"Grid_20160624000000000348_1"`
}

// MafraExaminResponse : mafra 농수축산 유통정보 조사가격(농수축산물표준코드변환)
type MafraExaminResponse struct {
	Grid201607220000000003521 struct {
		TotalCnt int `json:"totalCnt"`
		StartRow int `json:"startRow"`
		EndRow   int `json:"endRow"`
		Result   struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"result"`
		Row []struct {
			ROWNUM           int    `json:"ROW_NUM"`
			EXAMINDE         string `json:"EXAMIN_DE"`
			EXAMINSENM       string `json:"EXAMIN_SE_NM"`
			EXAMINSECODE     string `json:"EXAMIN_SE_CODE"`
			EXAMINAREANAME   string `json:"EXAMIN_AREA_NAME"`
			EXAMINAREACODE   string `json:"EXAMIN_AREA_CODE"`
			EXAMINMRKTNM     string `json:"EXAMIN_MRKT_NM"`
			EXAMINMRKTCODE   string `json:"EXAMIN_MRKT_CODE"`
			STDMRKTNM        string `json:"STD_MRKT_NM"`
			STDMRKTCODE      string `json:"STD_MRKT_CODE"`
			EXAMINPRDLSTNM   string `json:"EXAMIN_PRDLST_NM"`
			EXAMINPRDLSTCODE string `json:"EXAMIN_PRDLST_CODE"`
			EXAMINSPCIESNM   string `json:"EXAMIN_SPCIES_NM"`
			EXAMINSPCIESCODE string `json:"EXAMIN_SPCIES_CODE"`
			STDLCLASNM       string `json:"STD_LCLAS_NM"`
			STDLCLASCO       string `json:"STD_LCLAS_CO"`
			STDPRDLSTNM      string `json:"STD_PRDLST_NM"`
			STDPRDLSTCODE    string `json:"STD_PRDLST_CODE"`
			STDSPCIESNM      string `json:"STD_SPCIES_NM"`
			STDSPCIESCODE    string `json:"STD_SPCIES_CODE"`
			EXAMINUNITNM     string `json:"EXAMIN_UNIT_NM"`
			EXAMINUNIT       string `json:"EXAMIN_UNIT"`
			STDUNITNM        string `json:"STD_UNIT_NM"`
			STDUNITCODE      string `json:"STD_UNIT_CODE"`
			EXAMINGRADNM     string `json:"EXAMIN_GRAD_NM"`
			EXAMINGRADCODE   string `json:"EXAMIN_GRAD_CODE"`
			STDGRADNM        string `json:"STD_GRAD_NM"`
			STDGRADCODE      string `json:"STD_GRAD_CODE"`
			TODAYPRIC        int    `json:"TODAY_PRIC"`
			BFRTPRIC         int    `json:"BFRT_PRIC"`
			IMPTRADE         int    `json:"IMP_TRADE"`
			TRADEAMT         int    `json:"TRADE_AMT"`
		} `json:"row"`
	} `json:"Grid_20160722000000000352_1"`
}

// DataGoKrErrorResponse : 에러 발생 시 나오는 xml 구조체
type DataGoKrErrorResponse struct {
	XMLName      xml.Name `xml:"OpenAPI_ServiceResponse"`
	Text         string   `xml:",chardata"`
	CmmMsgHeader struct {
		Text             string `xml:",chardata"`
		ErrMsg           string `xml:"errMsg"`
		ReturnAuthMsg    string `xml:"returnAuthMsg"`
		ReturnReasonCode string `xml:"returnReasonCode"`
	} `xml:"cmmMsgHeader"`
}

// GetValueByColumnName : dataStruct  받아 컬럼값을 리턴
func GetValueByColumnName(dataStruct interface{}, columName string) string {
	r := reflect.ValueOf(dataStruct)
	value := reflect.Indirect(r).FieldByName(columName)

	return value.Interface().(string)
}

// IsRetryingDataGoKrResponse : 재시도 해야하는지 확인
// 정상 response이면 false 리턴
// 재시도 에러 response 이면 true 리턴
// 나머지는 프로그램 종료
func IsRetryingDataGoKrResponse(response string, metas string) bool {
	reasonCodesForRetrying := []string{"04", "01", "22", "500"}
	isError := strings.Contains(response, "errMsg")
	if isError {
		errorResponse := DataGoKrErrorResponse{}
		xml.Unmarshal([]byte(response), &errorResponse)

		reasonCode := errorResponse.CmmMsgHeader.ReturnReasonCode
		fmt.Println("reasonCode", reasonCode)
		fmt.Printf("%T\n", reasonCode)
		if util.InArray(reasonCode, reasonCodesForRetrying) {
			return true
		}

		errorMessage := fmt.Sprintf("%s, errorMesage: %s", metas, response)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(response))
	}

	return false
}

// RequestDataGoKr : DataGoKr API 요청
func RequestDataGoKr(url string, headers map[string]string, metas string) string {
	result := Request(url, headers)
	isRetrying := IsRetryingDataGoKrResponse(result, metas)

	numRetrying := config.Conf.NumRetrying
	var i int
	for i = 0; i < numRetrying; i++ {
		if !isRetrying {
			break
		}

		time.Sleep(time.Second * 360)
		result := Request(url, headers)
		isRetrying = IsRetryingDataGoKrResponse(result, metas)
		fmt.Printf("%d 번째 재시도\n", i+1)
	}
	if i == numRetrying {
		errorMessage := fmt.Sprintf("%s, 요청 실패", metas)
		alarm.PostMessage("default", errorMessage)
		panic(errors.New(errorMessage))
	}

	return result
}
