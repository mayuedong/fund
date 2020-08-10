package util

import (
	"fmt"
	"fund/entity"
	"io/ioutil"
	"math"
	"path"
	"strconv"
	"strings"
	"time"
)

type hisIndex struct {
	Time  string  `json:"time"`
	Close float64 `json:"close"`
	Price float64 `json:"price"`
}

var (
	g_setUpPlus   float64
	g_addPlus     float64
	g_mapHisIndex map[string][]*hisIndex
	g_hisIndex    *hisIndex
)

func GetHisIndex() *hisIndex {
	if nil == g_hisIndex {
		g_hisIndex = new(hisIndex)
		g_hisIndex.loadDb()
	}
	return g_hisIndex
}

func (this *hisIndex) loadDb() {
	t1 := time.Now().UnixNano()
	this.createTable()
	mapHisIndex := make(map[string][]*hisIndex)
	mysql := GetFundMysql()
	conf := entity.GetConf()
	files := conf.GetIndexData()
	for _, file := range files {
		rows, err := mysql.Db.Query(fmt.Sprintf(`SELECT * FROM %s WHERE "2010-12-31" < time AND time < "2020-01-01";`, file))
		if err != nil {
			entity.GetLog().Print(err)
			return
		}

		var sli []*hisIndex
		for rows.Next() {
			ptr := new(hisIndex)
			err := rows.Scan(&ptr.Time, &ptr.Close, &ptr.Price)
			if err != nil {
				entity.GetLog().Print(err)
				return
			}
			sli = append(sli, ptr)
		}
		mapHisIndex[file[1:]] = sli
		rows.Close()
	}
	g_mapHisIndex = mapHisIndex
	t2 := time.Now().UnixNano()
	entity.GetLog().Printf("count:%d cost:%d\n", len(mapHisIndex), (t2-t1)/1e6)
}

func (this *hisIndex) Update() {
	this.loadFile()
	this.insert()
	this.loadDb()
}

func (thiis *hisIndex) createTable() {
	mysql := GetFundMysql()
	conf := entity.GetConf()
	files := conf.GetIndexData()
	for _, name := range files {
		sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(`time` varchar(10), `close` float, `price` float, PRIMARY KEY(`time`))ENGINE=InnoDB DEFAULT CHARSET=utf8;", name)
		rows, err := mysql.Db.Query(sql)
		if err != nil {
			entity.GetLog().Print(sql, err)
		} else {
			rows.Close()
		}
	}
}
func (thiis *hisIndex) insert() {
	mysql := GetFundMysql()
	mapHisIndex := g_mapHisIndex
	for k, sli := range mapHisIndex {
		for _, v := range sli {
			sql := fmt.Sprintf(`REPLACE INTO %s(time,close,price)VALUES("%s",%f,%f);`, k, v.Time, v.Close, v.Price)
			rows, err := mysql.Db.Query(sql)
			if err != nil {
				entity.GetLog().Print(sql, err)
			} else {
				rows.Close()
			}
		}
		entity.GetLog().Printf(`table:%s count:%d`, k, len(sli))
	}
}

func (this *hisIndex) loadFile() {
	mapHisIndex := make(map[string][]*hisIndex)
	conf := entity.GetConf()
	files := conf.GetIndexData()
	for _, file := range files {
		filePath := path.Join("hisIndex", file)
		sliByte, err := ioutil.ReadFile(filePath)
		if nil != err {
			entity.GetLog().Print(err)
			return
		}
		str := string(sliByte)
		str = strings.TrimLeft(str, `["`)
		str = strings.TrimRight(str, `"]`)
		sliStr := strings.Split(str, `","`)
		for _, str := range sliStr {
			line := strings.Split(str, `,`)
			ptr := new(hisIndex)
			ptr.Time = line[0]
			ptr.Close, _ = strconv.ParseFloat(line[2], 64)
			ptr.Price, _ = strconv.ParseFloat(line[6], 64)
			ptr.Price /= 1e8

			sliIndex := mapHisIndex[file]
			if 0 == len(sliIndex) {
				var tmp []*hisIndex
				tmp = append(tmp, ptr)
				mapHisIndex[file] = tmp
			} else {
				sliIndex = append(sliIndex, ptr)
				mapHisIndex[file] = sliIndex
			}
		}
	}
	g_mapHisIndex = mapHisIndex
}

type (
	indexInfo struct {
		High   float64
		Medium float64
		Low    float64
		SeMe   float64
		Index  []float64
	}
)

var g_mapIndexInfo map[string]*indexInfo
var g_indexInfo *indexInfo

func GetMap() map[string]*indexInfo {
	return g_mapIndexInfo
}
func GetIndexInfo() *indexInfo {
	if nil == g_indexInfo {
		g_indexInfo = new(indexInfo)
		g_indexInfo.parse()
	}
	return g_indexInfo
}
func (this *indexInfo) GetHigh(k string) float64 {
	ptr := g_mapIndexInfo[k]
	if nil == ptr {
		return -1.0
	}
	return ptr.High
}
func (this *indexInfo) GetMedium(k string) float64 {
	ptr := g_mapIndexInfo[k]
	if nil == ptr {
		return -1.0
	}
	return ptr.Medium
}
func (this *indexInfo) GetSeMe(k string) float64 {
	ptr := g_mapIndexInfo[k]
	if nil == ptr {
		return -1.0
	}
	return ptr.SeMe
}
func (this *indexInfo) GetLow(k string) float64 {
	ptr := g_mapIndexInfo[k]
	if nil == ptr {
		return -1.0
	}
	return ptr.Low
}

func (this *indexInfo) Test() map[string]map[string]int {
	ret := make(map[string]map[string]int)
	conf := entity.GetConf()
	files := conf.GetIndexData()
	for _, file := range files {
		code := file[1:]
		if ptr, ok := g_mapIndexInfo[code]; ok {
			node := make(map[string]int)
			for i := 0; i < len(ptr.Index); i += len(ptr.Index) / 19 {
				cur := ptr.Index[i]
				node[strconv.Itoa(int(cur))] = this.GetRate(code, cur)
			}
			ret[code] = node
		}
	}
	return ret
}

func (this *indexInfo) GetRate(code string, cur float64) (price int) {
	ptr := g_mapIndexInfo[code]
	if nil == ptr {
		return -1
	}
	rate := 0.0
	if cur < ptr.SeMe {
		rate = (1.0-cur/ptr.SeMe)/(1.0-ptr.Low/ptr.SeMe)/g_addPlus + 1.0/g_setUpPlus
		return int(math.Pow(2.0, rate) * 100.0)

	} else if cur <= ptr.Medium {
		rate = (1.0 - cur/ptr.Medium) / (1.0 - ptr.SeMe/ptr.Medium) / g_setUpPlus
		return int(math.Pow(2.0, rate) * 100.0)
	}
	return 0
}

func (this *indexInfo) parse() {
	if 0 == len(g_mapHisIndex) {
		GetHisIndex()
	}
	mapIndexInfo := make(map[string]*indexInfo)
	for k, sli := range g_mapHisIndex {
		ptr := new(indexInfo)
		ptr.parseOne(k, sli)
		mapIndexInfo[k] = ptr
	}
	g_mapIndexInfo = mapIndexInfo
}

func (this *indexInfo) parseOne(k string, sli []*hisIndex) {
	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].Close > sli[j].Close {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	for _, v := range sli {
		this.Index = append(this.Index, v.Close)
	}
	l := len(this.Index)
	this.Low = this.Index[0]
	this.High = this.Index[l-1]
	this.Medium = this.Index[int(0.75*float64(l))]
	this.SeMe = this.High - (this.High-this.Low)*6.0/7.0
	if "000001" == k {
		g_setUpPlus = (200.0 / this.Medium) / (1.0 - this.SeMe/this.Medium)
		g_addPlus = (100.0 / this.SeMe) / (1.0 - this.Low/this.SeMe)
	}
}
