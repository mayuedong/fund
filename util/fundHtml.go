package util

import (
	"fund/entity"
	"golang.org/x/net/html"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type FundHtml struct {
	Id          string
	Name        string
	Scale       float64
	FoundDate   int
	ManagerDate int
	TrackAims   string
	TrackErr    float64
	Uptime      string
	Stock       map[string]string
}

func (r *FundHtml) Get() (sli []*FundHtml) {
	db := GetSqlite(fundHtmlTable)
	defer db.CLOSE()
	keys := db.KEYS()
	for _, key := range keys {
		ptr := new(FundHtml)
		db.GET(key, ptr)
		sli = append(sli, ptr)
	}
	return sli
}

func (r *FundHtml) get(id string) *FundHtml {
	db := GetSqlite(fundHtmlTable)
	defer db.CLOSE()
	ptr := new(FundHtml)
	ptr.Id = id
	db.GET(id, ptr)
	return ptr
}

func (this *FundHtml) getUptime() string {
	return this.Uptime
}

func (this *FundHtml) getWait() int {
	return -30
}

func (this *FundHtml) getUrl() string {
	strUrl := entity.GetConf().GetFundHtml()
	return strings.Replace(strUrl, `${code}`, this.Id, -1)
}

func (r *FundHtml) Update(ids []string) {
	var sli []APIUP
	for _, id := range ids {
		sli = append(sli, r.get(id))
	}
	setTask(sli)
}

func (r *FundHtml) GetScale(id string) float64 {
	ptr := r.get(id)
	if nil == ptr {
		return 0.0
	}
	return ptr.Scale
}

func (r *FundHtml) GetFoundDate(id string) int {
	ptr := r.get(id)
	if nil == ptr {
		return 0
	}
	return ptr.FoundDate
}

func (r *FundHtml) GetManagerDate(id string) int {
	ptr := r.get(id)
	if nil == ptr {
		return 0
	}
	return ptr.ManagerDate
}

func (r *FundHtml) GetTrackAims(id string) string {
	ptr := r.get(id)
	if nil == ptr {
		return ""
	}
	return ptr.TrackAims
}

func (r *FundHtml) GetStock(id string) map[string]string {
	ptr := r.get(id)
	if nil == ptr {
		return nil
	}
	return ptr.Stock
}

func (r *FundHtml) GetTrackErr(id string) float64 {
	ptr := r.get(id)
	if nil == ptr {
		return 99.9
	}
	return ptr.TrackErr
}

func (this *FundHtml) parse(b []byte) error {
	this.Name = this.parseName(string(b))
	this.Scale = this.parseScale(string(b))
	this.FoundDate = this.parseFoundDate(string(b))
	this.ManagerDate = this.parseManagerDate(string(b))
	this.TrackAims = this.parseTrackAims(string(b))
	this.TrackErr = this.parseTrackErr(string(b))
	this.Stock = this.parseStock(string(b))
	this.Uptime = time.Now().AddDate(0, 0, 5+rand.Intn(8)+rand.Intn(8)).String()[:len("2020-12-12")]
	db := GetSqlite(fundHtmlTable)
	defer db.CLOSE()
	return db.SET(this.Id, this)
}
func (r *FundHtml) parseName(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.StartTagToken == tt {
			t := z.Token()
			if "span" == t.Data {
				for _, a := range t.Attr {
					if "class" == a.Key && "fix_fname" == a.Val {
						for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
							if html.TextToken == tt {
								return strings.Trim(z.Token().Data, cutset)
							}
						}
					}
				}
			}
		}
	}
	return ""
}

func (r *FundHtml) parseScale(str string) float64 {
	str = r.parseScaleHtml(str)
	pos := strings.Index(str, "亿")
	if -1 == pos {
		entity.GetLog().Println("not found 亿 of")
		return 0.0
	}

	scale, err := strconv.ParseFloat(str[:pos], 10)
	if nil != err {
		entity.GetLog().Println("Parse Scale ", err)
		return 0.0
	}
	return scale
}

func (r *FundHtml) parseScaleHtml(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	trigger := false
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.TextToken == tt {
			data := z.Token().Data
			if trigger {
				return strings.Trim(data, cutset)
			}
			if strings.Contains(data, "基金规模") {
				trigger = true
			}
		}
	}
	return ""
}

func (this *FundHtml) parseManagerDate(str string) int {
	str = this.parseManagerDateHtml(str)
	pos := strings.Index(str, "~")
	if -1 == pos {
		entity.GetLog().Println("Parse managerDate", this.Id, this.Name)
		return 0
	}
	strDate := str[:pos]
	t, e := time.Parse("2006-01-02", strDate)
	if nil != e {
		entity.GetLog().Println("time parse ", e)
		return 0
	}
	return int(time.Since(t).Hours() / 24.0)
}
func (r *FundHtml) parseManagerDateHtml(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.StartTagToken == tt {
			t := z.Token()
			if "td" == t.Data {
				for _, a := range t.Attr {
					if "class" == a.Key && "td01" == a.Val {
						for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
							if html.TextToken == tt {
								return strings.Trim(z.Token().Data, cutset)
							}
						}
					}
				}
			}
		}
	}
	return ""
}

func (this *FundHtml) parseFoundDate(str string) int {
	str = this.parseFoundDateHtml(str)
	t, e := time.Parse("2006-01-02", str)
	if nil != e {
		entity.GetLog().Println("time parse", this.Id, this.Name, e)
		return 0
	}
	return int(time.Since(t).Hours() / 24.0)
}
func (r *FundHtml) parseFoundDateHtml(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	trigger := false
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.TextToken == tt {
			data := z.Token().Data
			if trigger {
				return strings.Trim(data, cutset)
			}
			if strings.Contains(data, "成 立 日") {
				trigger = true
			}
		}
	}
	return ""
}

func (r *FundHtml) parseTrackAims(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	trigger := false
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.TextToken == tt {
			data := z.Token().Data
			if trigger {
				return strings.Trim(data, cutset)
			}
			if strings.Contains(data, "跟踪标的") {
				trigger = true
			}
		}
	}
	return ""
}

func (this *FundHtml) parseTrackErr(str string) float64 {
	str = this.parseTrackErrHtml(str)
	strErr := strings.Trim(str, "%")
	fErr, err := strconv.ParseFloat(strErr, 10)
	if nil != err {
		entity.GetLog().Println("Parse trackErr ", this.Id, this.Name, err)
		return 99.9
	}
	return fErr
}
func (r *FundHtml) parseTrackErrHtml(str string) string {
	z := html.NewTokenizer(strings.NewReader(str))
	trigger := false
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		if html.TextToken == tt {
			data := z.Token().Data
			if trigger {
				return strings.Trim(data, cutset)
			}
			if strings.Contains(data, "跟踪误差") {
				trigger = true
			}
		}
	}
	return ""
}
func (r *FundHtml) parseStock(str string) map[string]string {
	m := make(map[string]string)
	z := html.NewTokenizer(strings.NewReader(str))
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		t := z.Token()
		if html.StartTagToken == tt && 2 == len(t.Attr) {
			str := ""
			for _, m := range t.Attr {
				str += m.Val
			}
			if "fund_item quotationItem_DataTable popTabquotationItem_DataTable" == str {
				for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
					t = z.Token()
					if html.StartTagToken == tt && 1 == len(t.Attr) {
						if "alignRight10" == t.Attr[0].Val {
							for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
								t = z.Token()
								if html.StartTagToken == tt && "tr" == t.Data {
									td := 0
									k := ""
									for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
										t = z.Token()
										if html.StartTagToken == tt && "td" == t.Data {
											td++
											if 4 == td {
												break
											}
											switch td {
											case 1:
												for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
													t = z.Token()
													if html.StartTagToken == tt && "a" == t.Data {
														tt = z.Next()
														if html.TextToken == tt {
															k = z.Token().Data
														}
													} else if html.EndTagToken == tt && "a" == t.Data {
														break
													}
												}
											case 2:
												tt = z.Next()
												if html.TextToken == tt {
													m[k] = z.Token().Data
												}
											}
										}
									}
								} else if html.EndTagToken == tt && "table" == t.Data {
									return m
								}
							}
						}
					}
				}
			}
		}
	}
	return m
}
