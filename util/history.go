package util

import (
	"fmt"
	"fund/entity"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type History struct {
}

func (r *History) Get() (sli []*Turnover) {
	db := GetSqlite(historyTable)
	defer db.CLOSE()
	keys := db.KEYS()
	for _, key := range keys {
		ptr := new(Turnover)
		db.GET(key, ptr)
		sli = append(sli, ptr)
	}
	return sli
}

func (r *History) Update(ids []string) {
	dir, _ := filepath.Abs("history")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		entity.GetLog().Fatal("not found path ", dir)
	}
	files, err := ioutil.ReadDir(dir)
	if nil != err {
		entity.GetLog().Print(err)
		return
	}

	for _, file := range files {
		fname := file.Name()
		filePath := path.Join(dir, fname)
		b, err := ioutil.ReadFile(filePath)
		if nil != err {
			entity.GetLog().Print(err)
			return
		}

		var sli []float64
		lines := strings.Split(string(b), "\n")
		for _, line := range lines {
			line = strings.TrimLeft(line, `",`)
			elements := strings.Split(line, ",")
			if 5 < len(elements) && "2010-12-31" < elements[0] {
				if price, err := strconv.ParseFloat(elements[2], 64); nil == err {
					sli = append(sli, price)
				}
			}
		}
		r.parse(fname, sli)
	}
}

func (r *History) parse(k string, sli []float64) {
	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i] > sli[j] {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}

	ptr := new(Turnover)
	ptr.Index = sli

	l := len(ptr.Index)
	ptr.Id = k
	ptr.Low = ptr.Index[0]
	ptr.High = ptr.Index[l-1]
	ptr.Medium = ptr.Index[int(0.75*float64(l))]
	ptr.SeMe = ptr.High - (ptr.High-ptr.Low)*6.0/7.0
	db := GetSqlite(historyTable)
	db.SET(k, ptr)
	db.CLOSE()
}

type Turnover struct {
	Id     string
	High   float64
	Medium float64
	Low    float64
	SeMe   float64
	Index  []float64
}

func (this *Turnover) Get(id string) *Turnover {
	db := GetSqlite(historyTable)
	defer db.CLOSE()
	ptr := new(Turnover)
	db.GET(id, ptr)
	return ptr
}

type TestInvest struct {
	Id       string            `json:"指数代码"`
	Name     string            `json:"指数名称"`
	Medium   float64           `json:"建仓点"`
	SeMe     float64           `json:"加仓点"`
	Interval int               `json:"间隔"`
	Invest   map[string]string `json:"建仓表"`
}

func (r *Turnover) Test() (ret []*TestInvest) {
	topics := entity.GetConf().GetIndexTopics()
	sli := new(History).Get()
	for _, ptr := range sli {
		node := new(TestInvest)
		node.Invest = make(map[string]string)
		node.Id = ptr.Id
		node.Medium = ptr.Medium
		node.SeMe = ptr.SeMe
		node.Name = topics[ptr.Id]
		interval := len(ptr.Index) / 40
		node.Interval = interval
		for i, j := 0, interval; j < len(ptr.Index); i, j = i+interval, j+interval {
			if ptr.GetRate(ptr.Index[i]) < int(entity.GetConf().GetBasePrice()) {
				break
			}
			key := fmt.Sprintf("%.1f~%.1f", ptr.Index[i], ptr.Index[j])
			val := fmt.Sprintf("%d~%d", ptr.GetRate(ptr.Index[i]), ptr.GetRate(ptr.Index[j]))
			node.Invest[key] = val
		}
		ret = append(ret, node)
	}
	return ret
}

func (this *Turnover) GetRate(cur float64) int {
	base := entity.GetConf().GetBasePrice()
	ptr := this.Get("000001")
	setUpPlus := (2.0 * base / ptr.Medium) / (1.0 - ptr.SeMe/ptr.Medium)
	addPlus := (1.0 * base / ptr.SeMe) / (1.0 - ptr.Low/ptr.SeMe)

	rate := 0.0
	if cur <= this.SeMe {
		rate = (1.0-cur/this.SeMe)/(1.0-this.Low/this.SeMe)/addPlus + 1.0/setUpPlus
	} else {
		rate = (1.0 - cur/this.Medium) / (1.0 - this.SeMe/this.Medium) / setUpPlus
	}
	return int(math.Pow(2.0, rate) * base)
}
