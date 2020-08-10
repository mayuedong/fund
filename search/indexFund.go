package search

import (
	"fund/entity"
	"fund/util"
	"strings"
)

type index struct {
	Id          string  `json:"基金代码"`
	Name        string  `json:"基金名称"`
	Scale       float64 `json:"基金规模"`
	HoldPe      float64 `json:"散户持有百分比"`
	FoundDate   int     `json:"基金成立天数"`
	TrackAims   string  `json:"跟踪标的"`
	TrackErr    float64 `json:"跟踪误差"`
	CostRate    float64 `json:"运作费用"`
	Manager     string  `json:"基金经理"`
	ManagerDate int     `json:"基金经理管理天数"`
	FundSize    string  `json:"管理规模"`
	Year1       float64 `json:"近一年收益"`
	Year2       float64 `json:"近二年收益"`
	Year3       float64 `json:"近三年收益"`
}

func (this *index) set(ptr *util.FundList) {
	this.Id = ptr.GetId()
	this.Name = ptr.GetName()
	this.Year1 = ptr.GetYear1()
	this.Year2 = ptr.GetYear2()
	this.Year3 = ptr.GetYear3()

	ptrMixInfo := util.GetMixInfo()
	if nil != ptrMixInfo {
		this.Manager = ptrMixInfo.GetName(this.Id)
		this.FundSize = ptrMixInfo.GetFundSize(this.Id)
		this.HoldPe = ptrMixInfo.GetHoldPe(this.Id)
	}
	ptrFundInfo := util.GetFundInfo()
	if nil != ptrFundInfo {
		this.Scale = ptrFundInfo.GetScale(this.Id)
		this.ManagerDate = ptrFundInfo.GetManagerDate(this.Id)
		this.TrackErr = ptrFundInfo.GetTrackErr(this.Id)
		this.TrackAims = ptrFundInfo.GetTrackAims(this.Id)
		this.FoundDate = ptrFundInfo.GetFoundDate(this.Id)
	}
	ptrCostRate := util.GetCostRate()
	if nil != ptrCostRate {
		this.CostRate = ptrCostRate.GetRate(this.Id)
	}
}

type indexPool struct {
	Pool map[string][]*index
}

func (this *indexPool) getIds() (ids []string) {
	for _, sli := range this.Pool {
		for _, v := range sli {
			ids = append(ids, v.Id)
		}
	}
	return ids
}

func (this *indexPool) filterScale(sli []*index) []*index {
	for i := 0; i < len(sli); i++ {
		if sli[i].Scale*sli[i].HoldPe*0.01 < 3.0 {
			sli[i] = sli[len(sli)-1]
			sli = sli[:len(sli)-1]
			i--
		}
	}
	return sli
}
func (this *indexPool) sortTrackErr(sli []*index) []*index {
	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			/*
				if sli[i].TrackErr > sli[j].TrackErr {
					sli[i], sli[j] = sli[j], sli[i]
				}
				if sli[i].TrackErr == sli[j].TrackErr && sli[i].Scale*sli[i].HoldPe < sli[j].Scale*sli[j].HoldPe {
					sli[i], sli[j] = sli[j], sli[i]
				}
			*/
			if sli[i].Scale*sli[i].HoldPe < sli[j].Scale*sli[j].HoldPe {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	return sli
}
func (this *indexPool) search(req *request) {
	this.Pool = make(map[string][]*index)
	ptrFund := util.GetFund()
	sliFund := ptrFund.GetIndex()
	conf := entity.GetConf()
	tips := conf.GetIndexTips()
	for _, v := range sliFund {
		ptrIndex := new(index)
		ptrIndex.set(v)

		for _, tip := range tips {
			if strings.Contains(ptrIndex.TrackAims, tip) {
				sliIndex := this.Pool[tip]
				sliIndex = append(sliIndex, ptrIndex)
				this.Pool[tip] = sliIndex
				break
			}
		}
	}
	for k, sli := range this.Pool {
		sli = this.filterScale(sli)
		this.Pool[k] = this.sortTrackErr(sli)
	}
}
