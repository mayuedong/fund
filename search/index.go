package search

import (
	"fund/entity"
	"fund/util"
	"strings"
)

type index struct {
	Id          string            `json:"基金代码"`
	Name        string            `json:"基金名称"`
	Scale       float64           `json:"基金规模"`
	HoldPe      float64           `json:"散户持有百分比"`
	FoundDate   int               `json:"基金成立天数"`
	TrackAims   string            `json:"跟踪标的"`
	TrackErr    float64           `json:"跟踪误差"`
	Rate        float64           `json:"运作费用"`
	Manager     string            `json:"基金经理"`
	ManagerDate int               `json:"基金经理管理天数"`
	FundSize    string            `json:"管理规模"`
	Year1       float64           `json:"近一年收益"`
	Year2       float64           `json:"近二年收益"`
	Year3       float64           `json:"近三年收益"`
	Stock       map[string]string `json:"重仓股"`
}

type indexPool struct {
	Pool      map[string][]*index
	Pool1     map[string][]*index
	TrackAims map[string][]string
}

func (this *indexPool) get() {
	this.Pool = make(map[string][]*index)
	this.Pool1 = make(map[string][]*index)
	this.TrackAims = make(map[string][]string)
	list := new(util.FundHtml).Get()
	fj := new(util.FundJs)
	li := new(util.IndexList)
	rate := new(util.Rate)
	tips := entity.GetConf().GetIndexTopics()
	for _, tip := range tips {
		var sli []*index
		for _, l := range list {
			if strings.Contains(l.TrackAims, tip) {
				tmp := new(index)
				tmp.Id = l.Id
				tmp.Name = l.Name
				tmp.Scale = l.GetScale(l.Id)
				tmp.TrackAims = l.GetTrackAims(l.Id)
				tmp.TrackErr = l.GetTrackErr(l.Id)
				tmp.FoundDate = l.GetFoundDate(l.Id)
				tmp.ManagerDate = l.GetManagerDate(l.Id)
				tmp.Stock = l.GetStock(l.Id)
				tmp.HoldPe = fj.GetHoldPe(l.Id)
				tmp.FundSize = fj.GetFundSize(l.Id)
				tmp.Manager = fj.GetManager(l.Id)
				tmp.Rate = rate.GetRate(l.Id)
				if ptr := li.GetOne(l.Id); nil != ptr {
					tmp.Year1 = ptr.Year1
					tmp.Year2 = ptr.Year2
					tmp.Year3 = ptr.Year3
				}
				sli = append(sli, tmp)
			}
		}
		this.Pool[tip] = sli
	}

	for _, l := range list {
		if 5 > l.Scale {
			continue
		}
		sli := this.TrackAims[l.TrackAims]
		sli = append(sli, l.Name)
		this.TrackAims[l.TrackAims] = sli
	}
	for k, sli := range this.TrackAims {
		if 3 > len(sli) || 1000 < len(sli) {
			delete(this.TrackAims, k)
		}
	}
}

func (this *indexPool) filter() {
	for key, sli := range this.Pool {
		var arr []*index
		for _, v := range sli {
			if v.HoldPe*v.Scale > 300.0 && !strings.Contains(v.Name, "增强") {
				arr = append(arr, v)
			}
		}
		this.Pool[key] = arr
	}
}

func (this *indexPool) sort() {
	for key, sli := range this.Pool {
		for i := 0; i < len(sli); i++ {
			for j := i + 1; j < len(sli); j++ {
				if sli[i].Scale < sli[j].Scale {
					sli[i], sli[j] = sli[j], sli[i]
				}
			}
		}
		this.Pool[key] = sli
		if 5 > len(sli) {
			this.Pool1[key] = this.Pool[key]
			delete(this.Pool, key)
			continue
		}
	}
}

func (this *indexPool) search(req *request) {
	this.get()
	this.filter()
	this.sort()
}
