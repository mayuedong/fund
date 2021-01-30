package search

import (
	//"fund/entity"
	"fund/util"
	//"strings"
)

type currency struct {
	Id          string  `json:"基金代码"`
	Name        string  `json:"基金名称"`
	Scale       float64 `json:"基金规模"`
	HoldPe      float64 `json:"散户持有百分比"`
	FoundDate   int     `json:"基金成立天数"`
	Manager     string  `json:"基金经理"`
	ManagerDate int     `json:"基金经理管理天数"`
	FundSize    string  `json:"管理规模"`
	Year1       float64 `json:"近一年收益"`
	Rate        float64 `json:"运作费用"`
	Price       float64 `json:"起购金额"`
}

type currencyPool struct {
	Pool []*currency `json:"货币基金"`
}

func (this *currencyPool) search(req *request) {
	this.get()
	this.filter()
	this.delSliNil()
	this.sort()
	this.Pool = this.Pool[:int(len(this.Pool)/3)]
}
func (this *currencyPool) get() {
	fh := new(util.FundHtml)
	fj := new(util.FundJs)
	list := new(util.CurrencyList).Get()
	rate := new(util.Rate)
	for _, v := range list {
		tmp := new(currency)
		tmp.Id = v.Id
		tmp.Name = v.Name
		tmp.Year1 = v.GetYear1()
		tmp.Price = v.GetPrice()
		tmp.Scale = fh.GetScale(v.Id)
		tmp.Rate = rate.GetRate(v.Id)
		tmp.FoundDate = fh.GetFoundDate(v.Id)
		tmp.ManagerDate = fh.GetManagerDate(v.Id)
		tmp.HoldPe = fj.GetHoldPe(v.Id)
		tmp.FundSize = fj.GetFundSize(v.Id)
		tmp.Manager = fj.GetManager(v.Id)
		this.Pool = append(this.Pool, tmp)
	}
}
func (this *currencyPool) filter() {
	for i, v := range this.Pool {
		if 10 > v.Scale || 1e5 < v.Price {
			this.Pool[i] = nil
		}
	}
}
func (this *currencyPool) delSliNil() {
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
func (this *currencyPool) sort() {
	for i := 0; i < len(this.Pool); i++ {
		for j := i + 1; j < len(this.Pool); j++ {
			if this.Pool[i].Year1 < this.Pool[j].Year1 {
				this.Pool[i], this.Pool[j] = this.Pool[j], this.Pool[i]
			}
		}
	}
}
