package entity

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var g_log *log.Logger

func GetLog() *log.Logger {
	if nil == g_log {
		g_log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return g_log
}

var g_client = &http.Client{}

func Get(url string) (ret []byte) {
	req, err := http.NewRequest("GET", url, nil)
	if nil != err {
		GetLog().Print(err)
		return
	}
	req.Header.Set("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.193 Safari/537.36`)

	resp, err := g_client.Do(req)
	if nil != err {
		GetLog().Print(err)
		return
	}

	defer resp.Body.Close()
	ret, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		GetLog().Print(err)
		return
	}
	return
}

type conf struct {
	IndexTips []string `json:"indexTips"`
	IndexData []string `json:"indexData"`
	IndexUrl  string   `json:"indexUrl"`
	MixList   string   `json:"mixList"`
	IndexList string   `json:"indexList"`
	FundInfo  string   `json:"fundInfo"`
	MixInfo   string   `json:"mixInfo"`
	CostRate  string   `json:"costRate"`
}

var g_conf *conf

func GetConf() *conf {
	return g_conf
}
func LoadConf(path string) error {
	g_conf = new(conf)
	return readFileUnMarshal(path, &g_conf)
}
func (this *conf) GetIndexTips() []string {
	return this.IndexTips
}
func (this *conf) GetIndexData() []string {
	return this.IndexData
}
func (this *conf) GetIndexUrl() string {
	return this.IndexUrl
}
func (this *conf) GetMixList() string {
	return this.MixList
}
func (this *conf) GetIndexList() string {
	return this.IndexList
}
func (this *conf) GetFundInfo() string {
	return this.FundInfo
}
func (this *conf) GetMixInfo() string {
	return this.MixInfo
}
func (this *conf) GetCostRate() string {
	return this.CostRate
}

func readFileUnMarshal(path string, ret interface{}) error {
	if b, e := ioutil.ReadFile(path); nil != e {
		return e
	} else {
		return json.Unmarshal(b, ret)
	}
}
