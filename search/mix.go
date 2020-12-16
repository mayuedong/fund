package search

import (
	"fund/util"
	"strconv"
	"strings"
)

type mix struct {
	Id          string            `json:"基金代码"`
	Name        string            `json:"基金名称"`
	Scale       float64           `json:"基金规模"`
	Manager     string            `json:"基金经理"`
	FundSize    string            `json:"管理规模"`
	ManagerDate int               `json:"基金经理管理天数"`
	FoundDate   int               `json:"基金成立天数"`
	HoldPe      float64           `json:"散户持有百分比"`
	Rate        float64           `json:"运作费用"`
	Year1       float64           `json:"近一年收益"`
	Year2       float64           `json:"近二年收益"`
	Year3       float64           `json:"近三年收益"`
	Stock       map[string]string `json:"重仓股"`
}

type mixPool struct {
	Pool           []*mix `json:"混合型基金"`
	Achievement    []*mix `json:"绩效优秀"`
	DisAchievement []*mix `json:"绩效稍逊"`
}

func (this *mixPool) search(req *request) {
	this.get()
	this.filter()
	this.delSliNil()
	this.sortOut()
	this.delSliNil()
}
func (this *mixPool) sortOut() {
	this.Achievement = this.filterAchievement(2)
	this.DisAchievement = this.disAchievement(2)
	var sli []*mix
	sli = append(sli, this.Achievement...)
	sli = append(sli, this.DisAchievement...)
	for _, s := range sli {
		for i, p := range this.Pool {
			if p == s {
				this.Pool[i] = nil
			}
		}
	}
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
	this.Pool = this.filterAchievement(4)
	this.filterScale()
	this.filterWorkTime()
	this.filterUserName()
	//this.filterFundSize()
	this.filterFoundTime()
}
func (this *mixPool) filterAchievement(rank int) []*mix {
	sli := this.Pool
	l := int(len(sli) / rank)
	var good1, good2, good3, good []*mix
	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].Year1 < sli[j].Year1 {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	good1 = append(good1, sli[:l]...)

	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].Year2 < sli[j].Year2 {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	good2 = append(good2, sli[:l]...)

	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].Year3 < sli[j].Year3 {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	good3 = append(good3, sli[:l]...)

	for _, g1 := range good1 {
		for _, g2 := range good2 {
			for _, g3 := range good3 {
				if g1 == g2 && g2 == g3 {
					good = append(good, g1)
				}
			}
		}
	}
	return good
}
func (this *mixPool) disAchievement(rank int) []*mix {
	sli := this.Pool
	l := int(len(sli) / rank)
	var good1, good2, good3, good []*mix
	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].Year1 < sli[j].Year1 {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	good1 = append(good1, sli[l:]...)

	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].Year2 < sli[j].Year2 {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	good2 = append(good2, sli[l:]...)

	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].Year3 < sli[j].Year3 {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	good3 = append(good3, sli[l:]...)

	for _, g1 := range good1 {
		for _, g2 := range good2 {
			for _, g3 := range good3 {
				if g1 == g2 && g2 == g3 {
					good = append(good, g1)
				}
			}
		}
	}
	return good
}
func (this *mixPool) filterScale() {
	for i, v := range this.Pool {
		if nil == v {
			continue
		}
		if 35.0 > v.Scale*0.01*v.HoldPe {
			this.Pool[i] = nil
		}
	}
}

func (this *mixPool) filterFoundTime() {
	for i, v := range this.Pool {
		if nil == v {
			continue
		}
		if v.FoundDate < 5*365 {
			this.Pool[i] = nil
		}
	}
}

func (this *mixPool) filterWorkTime() {
	for i, v := range this.Pool {
		if nil == v {
			continue
		}
		if v.ManagerDate < 2*365 {
			this.Pool[i] = nil
		}
	}
}

func (this *mixPool) filterFundSize() {
	for i, v := range this.Pool {
		if nil == v {
			continue
		}

		sliSize := strings.Split(v.FundSize, "(")
		if 1 < len(sliSize) {
			pos := strings.Index(sliSize[1], "只")
			if -1 != pos {
				num, _ := strconv.Atoi(sliSize[1][:pos])
				if 5 < num {
					this.Pool[i] = nil
				}
			}
		}
	}
}

func (this *mixPool) filterUserName() {
	for i, v := range this.Pool {
		if nil == v {
			continue
		}
		if strings.Contains(v.Manager, "&") {
			this.Pool[i] = nil
		}
	}
}

func (this *mixPool) delSliNil() {
	pos := -1
	for _, v := range this.Pool {
		if nil != v {
			pos++
			this.Pool[pos] = v
		}
	}
	this.Pool = append(this.Pool, nil)
	this.Pool = this.Pool[:pos+1]
}
