package util

import (
	"fmt"
	"fund/entity"
	"golang.org/x/net/html"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type CostRate struct {
	id       string
	costRate float64
	uptime   string
}

const (
	costRate = "costRate"
)

func (this *CostRate) createTable() {
	mysql := GetFundMysql()
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(`id` varchar(8) NOT NULL,  `costRate` float, `uptime` varchar(16), PRIMARY KEY(`id`))ENGINE=InnoDB DEFAULT CHARSET=utf8;", costRate)
	rows, err := mysql.Db.Query(sql)
	if err != nil {
		entity.GetLog().Fatal(sql, err)
	} else {
		rows.Close()
	}
}

var g_costRate *CostRate
var g_costRates map[string]*CostRate

func GetCostRate() *CostRate {
	if nil == g_costRate {
		g_costRate = new(CostRate)
		g_costRate.loadDb()
	}
	return g_costRate
}
func (this *CostRate) GetRate(id string) float64 {
	v := g_costRates[id]
	if nil == v {
		return 99.9
	}

	return v.costRate
}

func (this *CostRate) loadDb() {
	t1 := time.Now().UnixNano()
	this.createTable()
	mysql := GetFundMysql()
	rows, err := mysql.Db.Query(fmt.Sprintf("SELECT id,costRate,uptime FROM %s;", costRate))
	if err != nil {
		entity.GetLog().Fatal(err)
	}
	defer rows.Close()

	tmp := make(map[string]*CostRate)
	for rows.Next() {
		ptr := new(CostRate)
		err := rows.Scan(&ptr.id, &ptr.costRate, &ptr.uptime)
		if err != nil {
			entity.GetLog().Fatal(err)
			continue
		}
		tmp[ptr.id] = ptr
	}
	if 0 != len(tmp) {
		g_costRates = tmp
	}
	t2 := time.Now().UnixNano()
	entity.GetLog().Printf("costRate:%d cost:%d\n", len(tmp), (t2-t1)/1e6)
}

func (this *CostRate) insert() {
	mysql := GetFundMysql()
	sql := fmt.Sprintf(`REPLACE INTO %s(id,costRate,uptime)VALUES("%s","%f","%s");`, costRate, this.id, this.costRate, this.uptime)
	row, err := mysql.Db.Query(sql)
	if err != nil {
		entity.GetLog().Fatal(err)
	} else {
		row.Close()
	}
}
func (this *CostRate) Update(ids []string) {
	this.download(ids)
	this.loadDb()
}

func (this *CostRate) download(ids []string) {
	now := time.Now()
	prevMonth := now.AddDate(0, 0, g_intervalDay).String()[:len("2020-12-12")]
	strUrl := entity.GetConf().GetCostRate()
	for i, id := range ids {
		v := g_costRates[id]
		if nil != v && prevMonth < v.uptime {
			continue
		}
		time.Sleep(time.Duration(rand.Intn(5)+rand.Intn(5)+2) * time.Second)
		tmpUrl := strings.Replace(strUrl, `${code}`, id, -1)
		b := entity.Get(tmpUrl)
		rate := new(CostRate)
		rate.id = id
		rate.costRate = rate.parseCostRate(string(b))
		rate.uptime = now.String()[:len("2020-12-12")]
		rate.insert()
		if 0 == i%10 {
			entity.GetLog().Printf("all:%d down:%d", len(ids), i)
		}
	}
}

func (this *CostRate) parseCostRate(str string) float64 {
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
