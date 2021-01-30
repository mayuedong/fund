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

type (
	HisNode struct {
		Time  string  `json:"t"`
		Index float64 `json:"v"`
	}
	History struct {
		Id    string     `json:"id"`
		Nodes []*HisNode `json:"node"`
	}
)

func (r *History) Get() (ret []*History) {
	db := GetSqlite(historyTable)
	defer db.CLOSE()
	keys := db.KEYS()
	for _, key := range keys {
		p := new(History)
		db.GET(key, p)
		ret = append(ret, p)
	}
	return ret
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

		ph := new(History)
		ph.Id = fname
		lines := strings.Split(string(b), "\n")
		for _, line := range lines {
			line = strings.TrimLeft(line, `",`)
			elements := strings.Split(line, ",")
			if 5 < len(elements) && "2010-12-31" < elements[0] && elements[0] < "2021-01-01" {
				if price, err := strconv.ParseFloat(elements[2], 64); nil == err {
					node := new(HisNode)
					node.Time = elements[0]
					node.Index = price
					ph.Nodes = append(ph.Nodes, node)
				}
			}
		}
		db := GetSqlite(historyTable)
		db.SET(fname, ph)
		db.CLOSE()
	}
}

type HisInfo struct {
	History
	Asc   []*HisNode
	Max   int
	Min   int
	Build int
	Plus  int
}

var g_hisInfo []*HisInfo

func GetHisInfo() []*HisInfo {
	if 0 == len(g_hisInfo) {
		new(HisInfo).parse()
	}
	return g_hisInfo
}
func (r *HisInfo) parse() {
	sli := new(History).Get()
	for i := 0; i < len(sli); i++ {
		ptr := new(HisInfo)
		ptr.Id = sli[i].Id
		ptr.Nodes = sli[i].Nodes
		aux := make([]*HisNode, len(sli[i].Nodes))
		copy(aux, sli[i].Nodes)
		for i := 0; i < len(aux); i++ {
			for j := i + 1; j < len(aux); j++ {
				if aux[i].Index > aux[j].Index {
					aux[i], aux[j] = aux[j], aux[i]
				}
			}
		}
		ptr.Max = len(aux) - 1
		ptr.Min = 0
		ptr.Build = int(0.75 * float64(ptr.Max))
		ptr.Plus = int(0.3 * float64(ptr.Max))
		ptr.Asc = aux
		g_hisInfo = append(g_hisInfo, ptr)
	}
}

func (this *HisInfo) calcMultiple() (sumOut, sumIn, n float64) {
	max := this.Asc[this.Max].Index
	for i := this.Min; i <= this.Build; i++ {
		price := this.GetPrice(this.Asc[i].Index)
		sumIn += price
		sumOut += price * max / this.Asc[i].Index
		if price == math.NaN() {
			fmt.Println(this.Plus, this.Build)
		}
	}
	n = sumOut / sumIn
	return sumIn, sumOut, n
}

func (this *HisInfo) GetPrice(cur float64) float64 {
	build := this.Asc[this.Build].Index
	plus := this.Asc[this.Plus].Index
	min := this.Asc[this.Min].Index
	rate := 0.0
	if cur <= plus {
		rate = (plus - cur) / (plus - min)
		rate = rate*4.2 + 3.8
	} else {
		rate = (build - cur) / (build - plus)
		rate *= 3.8
	}
	return math.Pow(2.0, rate) * entity.GetConf().GetBasePrice()
}

func (r *HisInfo) Test() (ret []*Invest) {
	topics := entity.GetConf().GetIndexTopics()
	sli := GetHisInfo()
	for _, v := range sli {
		ptr := new(Invest)
		ret = append(ret, ptr)
		ptr.Id = v.Id
		ptr.Name = topics[v.Id]
		ptr.Min = v.Asc[v.Min].Index
		ptr.Max = v.Asc[v.Max].Index
		ptr.ExpectMultiple = ptr.Max / ptr.Min
		ptr.Build = v.Asc[v.Build].Index
		ptr.Plus = v.Asc[v.Plus].Index
		ptr.SumIn, ptr.SumOut, ptr.RealityMultiple = v.calcMultiple()
		ptr.MultipleErr = ptr.RealityMultiple / ptr.ExpectMultiple
		ptr.setGroup(v)
	}
	return ret
}

type Invest struct {
	Id              string   `json:"指数代码"`
	Name            string   `json:"指数名称"`
	Build           float64  `json:"建仓点"`
	Plus            float64  `json:"加仓点"`
	Max             float64  `json:"最高点"`
	Min             float64  `json:"最低点"`
	ExpectMultiple  float64  `json:"定额倍数"`
	SumIn           float64  `json:"本金"`
	SumOut          float64  `json:"本息"`
	RealityMultiple float64  `json:"差额倍数"`
	MultipleErr     float64  `json:"倍数误差"`
	Group           []string `json:"群落"`
}

func (this *Invest) setGroup(in *HisInfo) {
	for cur := in.Asc[in.Max].Index; cur >= in.Asc[in.Min].Index; cur -= 100 {
		price := in.GetPrice(cur)
		if entity.GetConf().GetBasePrice()*0.9 < price {
			this.Group = append(this.Group, fmt.Sprintf("%.1f:%.1f", cur, price))
		}
	}
}
