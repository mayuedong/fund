package search

import (
	"fund/util"
	"strconv"
	"strings"
)

type mix struct {
	Id          string  `json:"基金代码"`
	Name        string  `json:"基金名称"`
	Scale       float64 `json:"基金规模"`
	Manager     string  `json:"基金经理"`
	FundSize    string  `json:"管理规模"`
	ManagerDate int     `json:"基金经理管理天数"`
	FoundDate   int     `json:"基金成立天数"`
	HoldPe      float64 `json:"散户持有百分比"`
	Year1       float64 `json:"近一年收益"`
	Year2       float64 `json:"近二年收益"`
	Year3       float64 `json:"近三年收益"`
	rank1       float64
	rank2       float64
	rank3       float64
}

func (this *mix) isLtRankAnd(rank float64) bool {
	if this.rank1 <= rank && this.rank2 <= rank && this.rank3 <= rank {
		return true
	}
	return false
}

func (this *mix) set(ptr *util.FundList) {
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
	ptrInfo := util.GetFundInfo()
	if nil != ptrInfo {
		this.Scale = ptrInfo.GetScale(this.Id)
		this.FoundDate = ptrInfo.GetFoundDate(this.Id)
		this.ManagerDate = ptrInfo.GetManagerDate(this.Id)
	}
}

type mixPool struct {
	Pool       []*mix `json:"混合型基金"`
	RankGt20   []*mix `json:"排名大于20%"`
	WorkLt2    []*mix `json:"运作时间小于2年"`
	ManagerGt1 []*mix `json:"集体决策"`
	CountGt5   []*mix `json:"运作超5只基金"`
	ValidScale []*mix `json:"规模太大或太小"`
	HoldLt     []*mix `json:"机构持有大于25%"`
}

func (this *mixPool) getIds() (ids []string) {
	for _, k := range this.Pool {
		ids = append(ids, k.Id)
	}
	return ids
}
func (this *mixPool) search(req *request) {
	ptrFund := util.GetFund()
	sliFund := ptrFund.GetMix()
	for _, v := range sliFund {
		ptrMix := new(mix)
		ptrMix.set(v)
		this.Pool = append(this.Pool, ptrMix)
	}
	this.filterRank(20)
	this.filterWorkTime()
	this.filterUserName()
	this.filterFundSize()
	this.filterScale()
	//	this.filterFoundTime()
	this.filterHold()
	this.delSliNil()
	this.sortScale()
}

func (this *mixPool) sortScale() {
	for i := 0; i < len(this.Pool); i++ {
		for j := i + 1; j < len(this.Pool); j++ {
			if this.Pool[i].Scale > this.Pool[j].Scale {
				this.Pool[i], this.Pool[j] = this.Pool[j], this.Pool[i]
			}
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
			this.WorkLt2 = append(this.WorkLt2, this.Pool[i])
			this.Pool[i] = nil
		}
	}
}

func (this *mixPool) filterRank(pos float64) {
	this.setRank()
	var tmp []*mix
	for _, v := range this.Pool {
		if v.isLtRankAnd(pos) {
			tmp = append(tmp, v)
		} else {
			this.RankGt20 = append(this.RankGt20, v)
		}
	}
	this.Pool = tmp
}

func (this *mixPool) filterHold() {
	for i, v := range this.Pool {
		if nil == v {
			continue
		}
		if v.HoldPe < 75.0 {
			this.HoldLt = append(this.HoldLt, v)
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
					this.CountGt5 = append(this.CountGt5, this.Pool[i])
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
			this.ManagerGt1 = append(this.ManagerGt1, this.Pool[i])
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

func (this *mixPool) filterScale() {
	for i, v := range this.Pool {
		if nil == v {
			continue
		}
		if 35.0 > v.Scale || v.Scale > 120.0 {
			this.ValidScale = append(this.ValidScale, this.Pool[i])
			this.Pool[i] = nil

		}
	}
}

func (this *mixPool) setRank() {
	pool := this.Pool
	fl := float64(len(pool))
	for i := 0; i < len(pool); i++ {
		for j := i + 1; j < len(pool); j++ {
			if pool[i].Year1 < pool[j].Year1 {
				pool[i], pool[j] = pool[j], pool[i]
			}
		}
	}
	for i, ptr := range pool {
		ptr.rank1 = 100 * float64(i) / fl
	}

	for i := 0; i < len(pool); i++ {
		for j := i + 1; j < len(pool); j++ {
			if pool[i].Year2 < pool[j].Year2 {
				pool[i], pool[j] = pool[j], pool[i]
			}
		}
	}
	for i, ptr := range pool {
		ptr.rank2 = 100 * float64(i) / fl
	}

	for i := 0; i < len(pool); i++ {
		for j := i + 1; j < len(pool); j++ {
			if pool[i].Year3 < pool[j].Year3 {
				pool[i], pool[j] = pool[j], pool[i]
			}
		}
	}
	for i, ptr := range pool {
		ptr.rank3 = 100 * float64(i) / fl
	}
}
