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
	MixInfo struct {
		id       string
		name     string
		fundSize string
		uptime   string
		holdPe   float64
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

const (
	mixInfo = "mixInfo"
)

func (this *MixInfo) createTable() {
	mysql := GetFundMysql()
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(`id` varchar(8) NOT NULL, `name` varchar(64), `fundSize` varchar(64), `uptime` varchar(16), `holdpe` float, PRIMARY KEY(`id`))ENGINE=InnoDB DEFAULT CHARSET=utf8;", mixInfo)
	rows, err := mysql.Db.Query(sql)
	if err != nil {
		entity.GetLog().Fatal(sql, err)
	} else {
		rows.Close()
	}
}

func (this *MixInfo) GetFundSize(id string) string {
	v := g_MixInfos[id]
	if nil != v {
		return v.fundSize
	}
	return ""
}
func (this *MixInfo) GetName(id string) string {
	v := g_MixInfos[id]
	if nil != v {
		return v.name
	}
	return ""
}
func (this *MixInfo) GetHoldPe(id string) float64 {
	v := g_MixInfos[id]
	if nil != v {
		return v.holdPe
	}
	return 0.0
}

var g_MixInfo *MixInfo
var g_MixInfos map[string]*MixInfo
var g_intervalDay int = -365

func GetMixInfo() *MixInfo {
	if nil == g_MixInfo {
		g_MixInfo = new(MixInfo)
		g_MixInfo.loadDb()
	}
	return g_MixInfo
}

func (this *MixInfo) parseHold(str string) {
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
				this.holdPe = node.Data[len(node.Data)-1]
			}
		}
	}
}

func (this *MixInfo) parsePerson(str string) {
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
	for i, v := range sliPerson {
		if i == len(sliPerson)-1 {
			this.name += v.Name
			this.fundSize += v.FundSize
		} else {
			this.name += v.Name + "&"
			this.fundSize += v.FundSize + "&"
		}
	}
}

func (this *MixInfo) insert() {
	mysql := GetFundMysql()
	sql := fmt.Sprintf(`REPLACE INTO %s(id,name,fundSize,uptime,holdPe)VALUES("%s","%s","%s","%s",%f);`, mixInfo, this.id, this.name, this.fundSize, this.uptime, this.holdPe)
	fmt.Println(sql)
	rows, err := mysql.Db.Query(sql)
	if err != nil {
		entity.GetLog().Print(err)
	} else {
		rows.Close()
	}
}

func (this *MixInfo) loadDb() {
	t1 := time.Now().UnixNano()
	this.createTable()
	mysql := GetFundMysql()
	rows, err := mysql.Db.Query("SELECT * FROM MixInfo;")
	if err != nil {
		entity.GetLog().Fatal(err)
	}
	defer rows.Close()

	MixInfos := make(map[string]*MixInfo)
	for rows.Next() {
		ptr := new(MixInfo)
		err := rows.Scan(&ptr.id, &ptr.name, &ptr.fundSize, &ptr.uptime, &ptr.holdPe)
		if err != nil {
			entity.GetLog().Fatal(err)
		}
		MixInfos[ptr.id] = ptr
	}
	g_MixInfos = MixInfos

	t2 := time.Now().UnixNano()
	entity.GetLog().Printf("count:%d cost:%d\n", len(g_MixInfos), (t2-t1)/1e6)
}

func (this *MixInfo) download(ids []string) {
	now := time.Now()
	prevMonth := now.AddDate(0, 0, g_intervalDay).String()[:len("2020-12-12")]
	strUrl := entity.GetConf().GetMixInfo()
	for i, id := range ids {
		v := g_MixInfos[id]
		if nil != v && prevMonth < v.uptime {
			continue
		}
		time.Sleep(time.Duration(rand.Intn(5)+rand.Intn(5)+2) * time.Second)
		strNow := time.Now().String()[:len("2020-12-12 12:12:12")]
		strNow = strings.Replace(strNow, "-", "", -1)
		strNow = strings.Replace(strNow, ":", "", -1)
		tmpUrl := strings.Replace(strUrl, `${code}`, id, -1)
		tmpUrl += fmt.Sprintf("?v=%s", strNow)
		sliByte := entity.Get(tmpUrl)
		lines := strings.Split(string(sliByte), ";")
		ptrMixInfo := new(MixInfo)
		ptrMixInfo.id = id
		ptrMixInfo.uptime = now.String()[:len("2020-12-12")]
		ptrMixInfo.parseHold(lines[21])
		ptrMixInfo.parsePerson(lines[24])
		ptrMixInfo.insert()
		if 0 == i%10 {
			entity.GetLog().Printf("all:%d down:%d\n", len(ids), i)
		}
	}
}
func (this *MixInfo) Update(ids []string) {
	this.download(ids)
	this.loadDb()
}
