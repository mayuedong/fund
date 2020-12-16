package util

/*
import (
	"fund/entity"
	"golang.org/x/net/html"
	"net/url"
	"strconv"
	"strings"
)

type (
	Product struct {
		Id          string
		Name        string
		FoundDate   string
		FundType    string
		ProductType string
		TrackAims   string
		Scale       float64
	}
	IndexTopic struct {
		Id       string
		Name     string
		Products []*Product
	}
)

func (r *IndexTopic) Get() (sli []*IndexTopic) {
	db := GetSqlite(indexTopicTable)
	defer db.CLOSE()
	keys := db.KEYS()
	for _, key := range keys {
		ptr := new(IndexTopic)
		db.GET(key, ptr)
		sli = append(sli, ptr)
	}
	return sli

}

func (r *IndexTopic) Update(ids []string) {
	topics := entity.GetConf().GetIndexTopics()
	for k, v := range topics {
		ptr := new(IndexTopic)
		ptr.Id = k
		ptr.Name = v
		Download(ptr)
	}
}

func (this *IndexTopic) getUptime() string {
	return ""
}

func (this *IndexTopic) getWait() int {
	return -30
}

func (this *IndexTopic) getUrl() string {
	strUrl := entity.GetConf().GetIndexTopic()
	name := url.QueryEscape(this.Name)
	return strings.Replace(strUrl, `${code}`, name, -1)
}
func (this *IndexTopic) parse(b []byte) error {
	var sli []*Product
	str := string(b)
	z := html.NewTokenizer(strings.NewReader(str))
	for tt := z.Next(); html.ErrorToken != tt; tt = z.Next() {
		Data := z.Token().Data
		if html.StartTagToken == tt && "tbody" == Data {
			for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
				Data = z.Token().Data
				if html.StartTagToken == tt && "tr" == Data {
					i := 0
					ptr := new(Product)
					for tt = z.Next(); html.ErrorToken != tt; tt = z.Next() {
						Data = z.Token().Data
						if html.StartTagToken == tt && "td" == Data {
							tt = z.Next()
							if html.EndTagToken == tt {
								i++
							} else if html.TextToken == tt {
								i++
								switch i {
								case 1:
									ptr.Id = z.Token().Data
								case 2:
									ptr.Name = z.Token().Data
								case 3:
									ptr.FoundDate = z.Token().Data
								case 4:
									ptr.FundType = z.Token().Data
								case 5:
									ptr.ProductType = z.Token().Data
								case 6:
									ptr.TrackAims = z.Token().Data
								case 7:
									ptr.Scale, _ = strconv.ParseFloat(z.Token().Data, 64)
								}
							}
						} else if html.EndTagToken == tt && "tr" == Data {
							break
						}
					}
					sli = append(sli, ptr)
				} else if html.EndTagToken == tt && "tbody" == Data {
					if 0 == len(sli) {
						break
					}
					this.Products = sli
					db := GetSqlite(indexTopicTable)
					db.SET(this.Id, this)
					db.CLOSE()
					break
				}
			}
		}
	}
	return nil
}
*/
