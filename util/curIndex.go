package util

import (
	"encoding/json"
	"fund/entity"
	"strings"
	"time"
)

type (
	jsonCurIndex struct {
		Name     string  `json:"name"`
		Id       string  `json:"symbol"`
		Price    float64 `json:"price"`
		Turnover float64 `json:"turnover"`
		Volume   float64 `json:"volume"`
	}
	curIndex struct {
		sli map[string]*jsonCurIndex
	}
)

var g_curIndex *curIndex
var g_time time.Time

func GetCurIndex() *curIndex {
	if nil == g_curIndex {
		g_curIndex = new(curIndex)
		g_curIndex.download()
		g_time = time.Now()
	} else if 30*time.Second < time.Now().Sub(g_time) {
		g_curIndex.download()
		g_time = time.Now()
	}
	return g_curIndex
}

func (this *curIndex) GetTurnover(id string) float64 {
	for _, v := range this.sli {
		if v.Id == id {
			return v.Turnover
		}
	}
	return 0
}

func (this *curIndex) GetPrice(id string) float64 {
	for _, v := range this.sli {
		if v.Id == id {
			return v.Price
		}
	}
	return 0
}

func (this *curIndex) download() {
	conf := entity.GetConf()
	url := conf.GetIndexUrl()
	sliByte := entity.Get(url)
	strData := string(sliByte)
	sliData := strings.Split(strData, "(")
	if 1 < len(sliData) {
		strData = sliData[1]
	}
	sliData = strings.Split(strData, ")")
	if 0 < len(sliData) {
		strData = sliData[0]
	}
	err := json.Unmarshal([]byte(strData), &this.sli)
	if nil != err {
		entity.GetLog().Print(err)
	}
}
