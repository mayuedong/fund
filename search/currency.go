package search

import (
	"fund/util"
)

type currencyPool struct {
	Pool []fund `json:"货币基金"`
}

func (this *currencyPool) search(req *request) {
	this.get()
	this.Pool = filter(this.Pool, func(v fund) bool { return v.getScale() > 10 && v.getPrice() < 1e5 })
	this.Pool = sort(this.Pool, func(max, min fund) bool { return max.getYear1() > min.getYear1() })
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
