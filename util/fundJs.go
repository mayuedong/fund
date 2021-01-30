package util

import (
	"encoding/json"
	"fmt"
	"fund/entity"
	"math/rand"
	"strings"
	"time"
)

type (
	FundJs struct {
		Id       string
		Name     string
		FundSize string
		Uptime   string
		HoldPe   float64
	}
	holdNode struct {
		Name string    `json:"name"` // "机构持有比例"
		Data []float64 `json:"data"`
	}
	hold struct {
		Series []*holdNode `json:"series"`
	}

	person struct {
		Id       string `json:"id"`       //12
		Name     string `json:"name"`     //: "俞晓斌",
		FundSize string `json:"fundSize"` //: "83.82亿(18只基金)",
	}
)

func (r *FundJs) get(id string) *FundJs {
	db := GetSqlite(fundJsTable)
	defer db.CLOSE()
	ptr := new(FundJs)
	ptr.Id = id
	db.GET(id, ptr)
	return ptr
}
func (r *FundJs) GetManager(id string) string {
	ptr := r.get(id)
	if nil == ptr {
		return ""
	}
	return ptr.Name
}
func (r *FundJs) GetFundSize(id string) string {
	ptr := r.get(id)
	if nil == ptr {
		return ""
	}
	return ptr.FundSize
}
func (r *FundJs) GetHoldPe(id string) float64 {
	ptr := r.get(id)
	if nil == ptr {
		return 0
	}
	return ptr.HoldPe
}
func (this *FundJs) parseHold(str string) {
	sli := strings.Split(str, "=")
	if 2 != len(sli) {
		entity.GetLog().Print("持有人结构:", str)
		return
	}

	ptrHold := new(hold)
	err := json.Unmarshal([]byte(sli[1]), ptrHold)
	if nil != err {
		entity.GetLog().Print(err)
		return
	}

	series := ptrHold.Series
	for _, node := range series {
		if "个人持有比例" == node.Name {
			if 0 != len(node.Data) {
				this.HoldPe = node.Data[len(node.Data)-1]
			}
		}
	}
}

func (this *FundJs) parsePerson(str string) {
	sli := strings.Split(str, "=")
	if 2 != len(sli) {
		entity.GetLog().Print("基金经理:", str)
		return
	}
	var sliPerson []*person
	err := json.Unmarshal([]byte(sli[1]), &sliPerson)
	if nil != err {
		entity.GetLog().Print(err)
		return
	}
	this.Name = ""
	this.FundSize = ""
	for i, v := range sliPerson {
		if i == len(sliPerson)-1 {
			this.Name += v.Name
			this.FundSize += v.FundSize
		} else {
			this.Name += v.Name + "&"
			this.FundSize += v.FundSize + "&"
		}
	}
}

func (this *FundJs) parse(b []byte) error {
	lines := strings.Split(string(b), ";")
	this.Uptime = time.Now().AddDate(0, 0, 5+rand.Intn(8)+rand.Intn(8)).String()[:len("2020-12-12")]
	if 21 < len(lines) && len(lines) < 25 {
		this.parseHold(lines[19])
		this.parsePerson(lines[21])
	} else if 24 < len(lines) {
		this.parseHold(lines[21])
		this.parsePerson(lines[24])
	}
	db := GetSqlite(fundJsTable)
	defer db.CLOSE()
	return db.SET(this.Id, this)
}

func (r *FundJs) Update(ids []string) (sli []APIUP) {
	for _, id := range ids {
		sli = append(sli, r.get(id))
	}
	return sli
}
func (this *FundJs) getUptime() string {
	return this.Uptime
}

func (this *FundJs) getWait() int {
	return -30
}

func (this *FundJs) getUrl() string {
	strUrl := entity.GetConf().GetFundJs()
	strUrl = strings.Replace(strUrl, `${code}`, this.Id, -1)
	strNow := time.Now().String()[:len("2020-12-12 12:12:12")]
	strNow = strings.Replace(strNow, "-", "", -1)
	strNow = strings.Replace(strNow, ":", "", -1)
	strNow = strings.Replace(strNow, " ", "", -1)
	strUrl += fmt.Sprintf("v=%s", strNow)
	return strUrl
}
