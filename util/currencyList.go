package util

import (
	"encoding/json"
	"fund/entity"
	"strconv"
)

type (
	CurrencyList struct {
		Id    string `json:"FCODE"`
		Name  string `json:"SHORTNAME"`
		Year1 string `json:"SYL_1N"`
		Price string `json:"MINSG"`
	}
	CurrencyListData struct {
		Data []*CurrencyList `json:"Data"`
	}
)

func (r *CurrencyList) Get() (sli []*CurrencyList) {
	db := GetSqlite(currencyListTable)
	defer db.CLOSE()
	keys := db.KEYS()
	for _, key := range keys {
		ptr := new(CurrencyList)
		db.GET(key, ptr)
		sli = append(sli, ptr)
	}
	return sli
}

func (this *CurrencyList) GetPrice() float64 {
	price, _ := strconv.ParseFloat(this.Price, 64)
	return price
}
func (this *CurrencyList) GetYear1() float64 {
	year1, _ := strconv.ParseFloat(this.Year1, 64)
	return year1
}
func (r *CurrencyList) Update(ids []string) {
	Download(new(CurrencyList))
}
func (this *CurrencyList) getUptime() string {
	return ""
}

func (this *CurrencyList) getWait() int {
	return -30
}

func (this *CurrencyList) getUrl() string {
	return entity.GetConf().GetCurrencyList()
}

func (this *CurrencyList) parse(b []byte) error {
	tmp := new(CurrencyListData)
	if err := json.Unmarshal(b, tmp); nil != err {
		entity.GetLog().Println("download currency err : ", err)
		return err
	}
	db := GetSqlite(currencyListTable)
	defer db.CLOSE()
	for _, ptr := range tmp.Data {
		db.SET(ptr.Id, ptr)
	}
	return nil
}
