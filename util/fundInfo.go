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

type FundInfo struct {
	id          string
	name        string
	scale       string
	foundDate   string
	managerDate string
	trackAims   string
	trackErr    string
	uptime      string
}

const (
	fundInfo = "fundInfo"
)

func (this *FundInfo) createTable() {
	mysql := GetFundMysql()
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s`(`id` varchar(8) NOT NULL, `name` varchar(64), `scale` varchar(32), `foundDate` varchar(32), `managerDate` varchar(64),`trackAims` varchar(32),`trackErr` varchar(8),`uptime` varchar(16), PRIMARY KEY(`id`))ENGINE=InnoDB DEFAULT CHARSET=utf8;", fundInfo)
	rows, err := mysql.Db.Query(sql)
	if err != nil {
		entity.GetLog().Fatal(sql, err)
	} else {
		rows.Close()
	}
}

var g_fundInfo *FundInfo
var g_fundInfos map[string]*FundInfo
var g_cutset = ` |：`

func GetFundInfo() *FundInfo {
	if nil == g_fundInfo {
		g_fundInfo = new(FundInfo)
		g_fundInfo.loadDb()
	}
	return g_fundInfo
}
func (this *FundInfo) GetScale(id string) float64 {
	v := g_fundInfos[id]
	if nil == v {
		return 0.0
	}

	pos := strings.Index(v.scale, "亿")
	if -1 == pos {
		entity.GetLog().Println("not found 亿 of ", id)
		return 0.0
	}

	scale, err := strconv.ParseFloat(v.scale[:pos], 10)
	if nil != err {
		entity.GetLog().Println("Parse Scale ", v.name, err)
		return 0.0
	}
	return scale
}
func (this *FundInfo) GetFoundDate(id string) int {
	v := g_fundInfos[id]
	if nil == v {
		return 0
	}
	t, e := time.Parse("2006-01-02", v.foundDate)
	if nil != e {
		entity.GetLog().Println("time parse ", v.name, v.foundDate, e)
		return 0
	}
	return int(time.Since(t).Hours() / 24.0)
}
func (this *FundInfo) GetManagerDate(id string) int {
	v := g_fundInfos[id]
	if nil == v {
		return 0
	}
	pos := strings.Index(v.managerDate, "~")
	if -1 == pos {
		entity.GetLog().Println("Parse managerDate ", v.name)
		return 0
	}
	strDate := v.managerDate[:pos]
	t, e := time.Parse("2006-01-02", strDate)
	if nil != e {
		entity.GetLog().Println("time parse ", v.name, v.managerDate, e)
		return 0
	}
	return int(time.Since(t).Hours() / 24.0)
}
func (this *FundInfo) GetTrackAims(id string) string {
	if v, _ := g_fundInfos[id]; nil != v {
		return v.trackAims
	}
	return ""
}
func (this *FundInfo) GetTrackErr(id string) float64 {
	v := g_fundInfos[id]
	if nil == v {
		return 99.9
	}
	strErr := strings.Trim(v.trackErr, "%")
	fErr, err := strconv.ParseFloat(strErr, 10)
	if nil != err {
		entity.GetLog().Println("Parse trackErr ", v.name, err)
		return 99.9
	}
	return fErr
}

func (this *FundInfo) loadDb() {
	t1 := time.Now().UnixNano()
	this.createTable()
	mysql := GetFundMysql()
	rows, err := mysql.Db.Query(fmt.Sprintf("SELECT id,name,scale,foundDate,managerDate,trackAims,trackErr,uptime FROM %s;", fundInfo))
	if err != nil {
		entity.GetLog().Fatal(err)
	}
	defer rows.Close()

	fundInfos := make(map[string]*FundInfo)
	for rows.Next() {
		ptr := new(FundInfo)
		err := rows.Scan(&ptr.id, &ptr.name, &ptr.scale, &ptr.foundDate, &ptr.managerDate, &ptr.trackAims, &ptr.trackErr, &ptr.uptime)
		if err != nil {
			entity.GetLog().Fatal(err)
			continue
		}
		fundInfos[ptr.id] = ptr
	}
	if 0 != len(fundInfos) {
		g_fundInfos = fundInfos
	}
	t2 := time.Now().UnixNano()
	entity.GetLog().Printf("fundInfos:%d cost:%d\n", len(fundInfos), (t2-t1)/1e6)
}

func (this *FundInfo) insert() {
	mysql := GetFundMysql()
	sql := fmt.Sprintf(`REPLACE INTO %s(id,name,scale,foundDate,managerDate,trackAims,trackErr,uptime)VALUES("%s","%s","%s","%s","%s","%s","%s","%s");`, fundInfo, this.id, this.name, this.scale, this.foundDate, this.managerDate, this.trackAims, this.trackErr, this.uptime)
	fmt.Println(sql)
	row, err := mysql.Db.Query(sql)
	if err != nil {
		entity.GetLog().Fatal(err)
	} else {
		row.Close()
	}
}
func (this *FundInfo) Update(ids []string) {
	this.download(ids)
	this.loadDb()
}

func (this *FundInfo) download(ids []string) {
	now := time.Now()
	prevMonth := now.AddDate(0, 0, g_intervalDay).String()[:len("2020-12-12")]
	strUrl := entity.GetConf().GetFundInfo()
	for i, id := range ids {
		v := g_fundInfos[id]
		if nil != v && prevMonth < v.uptime {
			continue
		}

		time.Sleep(time.Duration(rand.Intn(5)+rand.Intn(5)+2) * time.Second)
		tmpUrl := strings.Replace(strUrl, `${code}`, id, -1)
		b := entity.Get(tmpUrl)
		info := new(FundInfo)
		info.name = info.parseName(string(b))
		info.id = info.parseCode(string(b))
		info.scale = info.parseScale(string(b))
		info.foundDate = info.parseFoundDate(string(b))
		info.managerDate = info.parseManagerDate(string(b))
		info.trackAims = info.parseTrackAims(string(b))
		info.trackErr = info.parseTrackErr(string(b))
		info.uptime = now.String()[:len("2020-12-12")]
		info.insert()
		if 0 == i%10 {
			entity.GetLog().Printf("all:%d down:%d", len(ids), i)
		}
	}
}

func (this *FundInfo) parseName(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.StartTagToken == tt {
			t := z.Token()
			if "span" == t.Data {
				for _, a := range t.Attr {
					if "class" == a.Key && "fix_fname" == a.Val {
						for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
							if html.TextToken == tt {
								return strings.Trim(z.Token().Data, g_cutset)
							}
						}
					}
				}
			}
		}
	}
	return ""
}

func (this *FundInfo) parseCode(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.StartTagToken == tt {
			t := z.Token()
			if "span" == t.Data {
				for _, a := range t.Attr {
					if "class" == a.Key && "fix_fcode" == a.Val {
						for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
							if html.TextToken == tt {
								return strings.Trim(z.Token().Data, g_cutset)
							}
						}
					}
				}
			}
		}
	}
	return ""
}

func (this *FundInfo) parseScale(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	trigger := false
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.TextToken == tt {
			data := z.Token().Data
			if trigger {
				return strings.Trim(data, g_cutset)
			}
			if strings.Contains(data, "基金规模") {
				trigger = true
			}
		}
	}
	return ""
}

func (this *FundInfo) parseManagerDate(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.StartTagToken == tt {
			t := z.Token()
			if "td" == t.Data {
				for _, a := range t.Attr {
					if "class" == a.Key && "td01" == a.Val {
						for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
							if html.TextToken == tt {
								return strings.Trim(z.Token().Data, g_cutset)
							}
						}
					}
				}
			}
		}
	}
	return ""
}

func (this *FundInfo) parseFoundDate(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	trigger := false
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.TextToken == tt {
			data := z.Token().Data
			if trigger {
				return strings.Trim(data, g_cutset)
			}
			if strings.Contains(data, "成 立 日") {
				trigger = true
			}
		}
	}
	return ""
}

func (this *FundInfo) parseTrackAims(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	trigger := false
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.TextToken == tt {
			data := z.Token().Data
			if trigger {
				return strings.Trim(data, g_cutset)
			}
			if strings.Contains(data, "跟踪标的") {
				trigger = true
			}
		}
	}
	return ""
}
func (this *FundInfo) parseTrackErr(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	trigger := false
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.TextToken == tt {
			data := z.Token().Data
			if trigger {
				return strings.Trim(data, g_cutset)
			}
			if strings.Contains(data, "跟踪误差") {
				trigger = true
			}
		}
	}
	return ""
}
