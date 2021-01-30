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
	IndexTopics  map[string]string `json:"indexTopics"`
	CurIndex     string            `json:"curIndex"`
	IndexTopic   string            `json:"indexTopic"`
	MixList      string            `json:"mixList"`
	IndexList    string            `json:"indexList"`
	CurrencyList string            `json:"currencyList"`
	FundHtml     string            `json:"fundHtml"`
	FundJs       string            `json:"fundJs"`
	RateUrl      string            `json:"rateUrl"`
	BasePrice    float64           `json:"basePrice"`
	Rate         float64           `json:"rate"`
	Limit        float64           `json:"limit"`
	ForceUpdate  bool              `json:"forceUpdate"`
}

var g_conf *conf

func GetConf() *conf {
	return g_conf
}
func LoadConf(path string) error {
	g_conf = new(conf)
	return readFileUnMarshal(path, &g_conf)
}
func (this *conf) GetRate() float64 {
	return this.Rate
}
func (this *conf) GetLimit() float64 {
	return this.Limit
}
func (this *conf) GetIndexTopics() map[string]string {
	return this.IndexTopics
}
func (this *conf) GetCurIndex() string {
	return this.CurIndex
}
func (this *conf) GetBasePrice() float64 {
	return this.BasePrice
}
func (this *conf) GetForceUpdate() bool {
	return this.ForceUpdate
}
func (this *conf) GetIndexTopic() string {
	return this.IndexTopic
}
func (this *conf) GetMixList() string {
	return this.MixList
}
func (this *conf) GetCurrencyList() string {
	return this.CurrencyList
}
func (this *conf) GetIndexList() string {
	return this.IndexList
}
func (this *conf) GetFundHtml() string {
	return this.FundHtml
}
func (this *conf) GetFundJs() string {
	return this.FundJs
}
func (this *conf) GetRateUrl() string {
	return this.RateUrl
}

func readFileUnMarshal(path string, ret interface{}) error {
	if b, e := ioutil.ReadFile(path); nil != e {
		return e
	} else {
		return json.Unmarshal(b, ret)
	}
}
