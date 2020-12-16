package util

import (
	"fund/entity"
	"golang.org/x/net/html"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Rate struct {
	Id     string
	Rate   float64
	Uptime string
}

func (r *Rate) get(id string) *Rate {
	db := GetSqlite(rateTable)
	defer db.CLOSE()
	ptr := new(Rate)
	ptr.Id = id
	db.GET(id, ptr)
	return ptr
}

func (r *Rate) GetRate(id string) float64 {
	ptr := r.get(id)
	if nil == ptr {
		return 99.9
	}
	return ptr.Rate
}

func (r *Rate) Update(ids []string) {
	var sli []APIUP
	for _, id := range ids {
		sli = append(sli, r.get(id))
	}
	setTask(sli)
}

func (this *Rate) getUptime() string {
	return this.Uptime
}

func (this *Rate) getWait() int {
	return -30
}

func (this *Rate) getUrl() string {
	strUrl := entity.GetConf().GetRate()
	return strings.Replace(strUrl, `${code}`, this.Id, -1)
}

func (this *Rate) parse(b []byte) error {
	this.Rate = this.parseCostRate(string(b))
	this.Uptime = time.Now().AddDate(0, 0, 5+rand.Intn(8)+rand.Intn(8)).String()[:len("2020-12-12")]
	db := GetSqlite(rateTable)
	defer db.CLOSE()
	return db.SET(this.Id, this)
}

func (r *Rate) parseCostRate(str string) float64 {
	sum := 0
	sumRate := 0.0
	trigger := false
	z := html.NewTokenizer(strings.NewReader(str))
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.TextToken == tt {
			data := z.Token().Data
			if trigger {
				trigger = false
				strRate := strings.Trim(data, `%（每年）`)
				fRate, err := strconv.ParseFloat(strRate, 10)
				if nil != err {
					entity.GetLog().Printf("parse costRate:", err)
					continue
				}
				sumRate += fRate
				if 3 == sum {
					break
				}
			}
			if strings.Contains(data, "管理费率") {
				trigger = true
				sum++
			} else if strings.Contains(data, "托管费率") {
				trigger = true
				sum++
			} else if strings.Contains(data, "销售服务费率") {
				trigger = true
				sum++
			}
		}
	}
	if 0 != sumRate {
		return sumRate
	}
	return 99.9
}
