package util

import (
	"fmt"
	"fund/entity"
	"strconv"
	"strings"
	"time"
)

type FundList struct {
	id    string
	name  string
	year1 float64
	year2 float64
	year3 float64
}

const (
	mixList   = "mixList"
	indexList = "indexList"
)

func (this *FundList) createTable(name string) {
	mysql := GetFundMysql()
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(`id` varchar(8) NOT NULL, `name` varchar(64), `year1` float, `year2` float, `year3` float, PRIMARY KEY(`id`))ENGINE=InnoDB DEFAULT CHARSET=utf8;", name)
	rows, err := mysql.Db.Query(sql)
	if err != nil {
		entity.GetLog().Print(sql, err)
	} else {
		rows.Close()
	}
}

func (this *FundList) GetId() string {
	return this.id
}
func (this *FundList) GetName() string {
	return this.name
}
func (this *FundList) GetYear1() float64 {
	return this.year1
}
func (this *FundList) GetYear2() float64 {
	return this.year2
}
func (this *FundList) GetYear3() float64 {
	return this.year3
}

var g_fund *FundList
var g_indexFund, g_mixFund []*FundList

func GetFund() *FundList {
	if nil == g_fund {
		g_fund = new(FundList)
		g_fund.loadDb(mixList)
		g_fund.loadDb(indexList)
	}
	return g_fund
}

func (this *FundList) GetIndex() []*FundList {
	return g_indexFund
}
func (this *FundList) GetMix() []*FundList {
	return g_mixFund
}

func (ptrFund *FundList) set(sli []string) bool {
	if 13 > len(sli) {
		return false
	}
	ptrFund.id = sli[0]
	ptrFund.name = sli[1]
	ptrFund.year1, _ = strconv.ParseFloat(sli[10], 64)
	ptrFund.year2, _ = strconv.ParseFloat(sli[11], 64)
	ptrFund.year3, _ = strconv.ParseFloat(sli[12], 64)
	return true
}

func (this *FundList) insert(tableName string) {
	mysql := GetFundMysql()
	sql := fmt.Sprintf(`REPLACE INTO %s(name,year1,year2,year3,id)VALUES("%s",%f,%f,%f,"%s");`, tableName, this.name, this.year1, this.year2, this.year3, this.id)
	rows, err := mysql.Db.Query(sql)
	if err != nil {
		entity.GetLog().Print(err)
	} else {
		rows.Close()
	}
}

func (this *FundList) loadDb(tableName string) {
	t1 := time.Now().UnixNano()
	this.createTable(tableName)
	mysql := GetFundMysql()
	rows, err := mysql.Db.Query(fmt.Sprintf("SELECT id,name,year1,year2,year3 FROM %s;", tableName))
	if err != nil {
		entity.GetLog().Print(err)
	}
	defer rows.Close()

	var mix, index []*FundList
	for rows.Next() {
		ptr := new(FundList)
		err := rows.Scan(&ptr.id, &ptr.name, &ptr.year1, &ptr.year2, &ptr.year3)
		if err != nil {
			entity.GetLog().Print(err)
			continue
		}
		if indexList == tableName {
			index = append(index, ptr)
		} else if mixList == tableName {
			mix = append(mix, ptr)
		}
	}
	if 0 != len(mix) {
		g_mixFund = mix
	}
	if 0 != len(index) {
		g_indexFund = index
	}

	t2 := time.Now().UnixNano()
	entity.GetLog().Printf("index:%d mix:%d cost:%d\n", len(index), len(mix), (t2-t1)/1e6)
}

func (this *FundList) downIndex() {
	sliByte := entity.Get(entity.GetConf().GetIndexList())
	sliFund := parse(sliByte, `["`, `"]`, `","`)
	if 0 == len(sliFund) {
		entity.GetLog().Print("pase index ", string(sliByte))
	}
	for _, v := range sliFund {
		sli := strings.Split(v, "|")
		ptrFund := new(FundList)
		if !ptrFund.set(sli) {
			entity.GetLog().Println("parse fundList err : ", v)
			continue
		}
		ptrFund.insert(indexList)
	}
}

func (this *FundList) downMix() {
	sliByte := entity.Get(entity.GetConf().GetMixList())
	sliFund := parse(sliByte, `["`, `"]`, `","`)
	if 0 == len(sliFund) {
		entity.GetLog().Print("pase mix ", string(sliByte))
	}
	for _, v := range sliFund {
		sli := strings.Split(v, "|")
		ptrFund := new(FundList)
		if !ptrFund.set(sli) {
			entity.GetLog().Print("parse fundList err : ", v)
			continue
		}
		ptrFund.insert(mixList)
	}
}

func (this *FundList) Update() {
	this.downMix()
	this.downIndex()
	this.loadDb(mixList)
	this.loadDb(indexList)
}

func parse(strData []byte, start, end, sep string) (sli []string) {
	str := string(strData)
	pos := strings.Index(str, start)
	if -1 != pos {
		str = str[pos+len(start):]
	}
	pos = strings.Index(str, end)
	if -1 != pos {
		str = str[:pos]
	}
	return strings.Split(str, sep)
}
