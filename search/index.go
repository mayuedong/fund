package search

import (
	"fund/util"
	"strings"
)

type indexPool struct {
	Pool map[string][]fund
}

func (this *indexPool) get() {
	this.Pool = make(map[string][]fund)
	list := new(util.FundHtml).Get()
	fj := new(util.FundJs)
	li := new(util.IndexList)
	rate := new(util.Rate)

	for _, l := range list {
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
		sli := this.Pool[tmp.TrackAims]
		sli = append(sli, tmp)
		this.Pool[tmp.TrackAims] = sli
	}
}

func (this *indexPool) filter() {
	for key, sli := range this.Pool {
		sli = filter(sli, func(v fund) bool {
			return (!strings.Contains(v.getTrackAims(), "债")) && (!strings.Contains(v.getName(), "增强")) && (v.getScale() > 5.0)
		})
		if 3 > len(sli) || 1000 < len(sli) {
			delete(this.Pool, key)
		} else {
			this.Pool[key] = sli
		}
	}
}
func (this *indexPool) sort() {
	for key, sli := range this.Pool {
		sli = sort(sli, func(max, min fund) bool { return max.getScale() > min.getScale() })
		this.Pool[key] = sli
	}
}

func (this *indexPool) search(req *request) {
	this.get()
	this.filter()
	this.sort()
}
