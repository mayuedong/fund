package search

import (
	"fund/util"
	"strconv"
	"strings"
)

type mixPool struct {
	Pool    []fund `json:"混合型基金"`
	PoolMin []fund `json:"小盘"`
	PoolMax []fund `json:"大盘"`
}

func (this *mixPool) search(req *request) {
	this.get()
	this.filter()
}
func (this *mixPool) get() {
	fh := new(util.FundHtml)
	fj := new(util.FundJs)
	list := new(util.MixList).Get()
	rate := new(util.Rate)
	for _, v := range list {
		tmp := new(mix)
		tmp.Id = v.Id
		tmp.Name = v.Name
		tmp.Year1 = v.Year1
		tmp.Year2 = v.Year2
		tmp.Year3 = v.Year3
		tmp.Rate = rate.GetRate(v.Id)
		tmp.UpTime = fh.GetUpTime(v.Id)
		tmp.Scale = fh.GetScale(v.Id)
		tmp.FoundDate = fh.GetFoundDate(v.Id)
		tmp.ManagerDate = fh.GetManagerDate(v.Id)
		tmp.Stock = fh.GetStock(v.Id)
		tmp.HoldPe = fj.GetHoldPe(v.Id)
		tmp.FundSize = fj.GetFundSize(v.Id)
		tmp.Manager = fj.GetManager(v.Id)
		this.Pool = append(this.Pool, tmp)
	}
}
func (this *mixPool) filter() {
	this.Pool = this.filterAchievement(2)
	this.Pool = filter(this.Pool, func(v fund) bool { return 40.0 < v.getScale() })
	this.Pool = filter(this.Pool, func(v fund) bool { return v.getManagerDate() > 2*365 })
	this.Pool = filter(this.Pool, func(v fund) bool { return !strings.Contains(v.getManager(), "&") })
	this.Pool = filter(this.Pool, func(v fund) bool { return v.getFoundDate() > 5*365 })
	this.Pool = filter(this.Pool, func(v fund) bool {
		sliSize := strings.Split(v.getFundSize(), "(")
		if 1 < len(sliSize) {
			pos := strings.Index(sliSize[1], "只")
			if -1 != pos {
				num, _ := strconv.Atoi(sliSize[1][:pos])
				return num < 7
			}
		}
		return true
	})
	this.Pool = this.filterAchievement(2)
	this.PoolMax = filter(this.Pool, func(v fund) bool { return v.getScale() >= 100 })
	this.PoolMin = filter(this.Pool, func(v fund) bool { return v.getScale() < 100 })
	this.Pool = nil
}
func (this *mixPool) filterAchievement(rank int) []fund {
	good1 := sort(this.Pool, func(max, min fund) bool { return max.getYear1() > min.getYear1() })
	good2 := sort(this.Pool, func(max, min fund) bool { return max.getYear2() > min.getYear2() })
	good3 := sort(this.Pool, func(max, min fund) bool { return max.getYear3() > min.getYear3() })

	l := int(len(this.Pool)/rank + 1)
	good1 = good1[:l]
	good2 = good2[:l]
	good3 = good3[:l]

	var good []fund
	for _, g3 := range good3 {
	LOOP:
		for _, g2 := range good2 {
			if g2 == g3 {
				for _, g1 := range good1 {
					if g1 == g2 {
						good = append(good, g3)
						break LOOP
					}
				}
			}
		}
	}
	return good
}

func (this *mixPool) delSliNil() {
	pos := 0
	for _, v := range this.Pool {
		if nil != v {
			this.Pool[pos] = v
			pos++
		}
	}
	this.Pool = this.Pool[:pos]
}
