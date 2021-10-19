package models

// Kamisitem 카미스 아이템 데이터 모델
type Kamisitem struct {
	itemname    string
	kindname    string
	countryname string
	marketname  string
	yyyy        string
	regday      string
	price       string
}

type KamisApi struct {
	Condition []Condition `json:"condition"`
	Data      Data        `json:"data"`
}

type Condition struct {
	PStartday         string        `json:"p_startday"`
	PEndday           string        `json:"p_endday"`
	PItemcategorycode string        `json:"p_itemcategorycode"`
	PItemcode         string        `json:"p_itemcode"`
	PKindcode         string        `json:"p_kindcode"`
	PProductrankcode  string        `json:"p_productrankcode"`
	PCountycode       []interface{} `json:"p_countycode"`
	PConvertKgYn      string        `json:"p_convert_kg_yn"`
	PKey              string        `json:"p_key"`
	PID               string        `json:"p_id"`
	PReturntype       string        `json:"p_returntype"`
}

type Item struct {
	Itemname   string `json:"itemname"`
	Kindname   string `json:"kindname"`
	Countyname string `json:"countyname"`
	Marketname string `json:"marketname"`
	Yyyy       string `json:"yyyy"`
	Regday     string `json:"regday"`
	Price      string `json:"price"`
}

type Data struct {
	ErrorCode string `json:"error_code"`
	Item      []Item `json:"item"`
}
