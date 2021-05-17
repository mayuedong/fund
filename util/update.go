package util

import (
	"fund/entity"
)

type Update struct {
	mix      []*MixList
	index    []*IndexList
	currency []*CurrencyList
	history  []*History
}

func (this *Update) getAllId() (ids []string) {
	for _, v := range this.index {
		ids = append(ids, v.Id)
	}
	for _, v := range this.mix {
		ids = append(ids, v.Id)
	}
	for _, v := range this.currency {
		ids = append(ids, v.Id)
	}
	return ids
}

func (this *Update) AutoUp() {
	history := new(History)
	this.history = history.Get()
	if 0 == len(this.history) {
		history.Update(nil)
		this.history = history.Get()
		if 0 == len(this.history) {
			entity.GetLog().Fatal("Update history error")
		}
	}

	mixList := new(MixList)
	this.mix = mixList.Get()
	if 0 == len(this.mix) {
		mixList.Update(nil)
		this.mix = mixList.Get()
		if 0 == len(this.mix) {
			entity.GetLog().Fatal("Update mix error")
		}
	}

	indexList := new(IndexList)
	this.index = indexList.Get()
	if 0 == len(this.index) {
		indexList.Update(nil)
		this.index = indexList.Get()
		if 0 == len(this.index) {
			entity.GetLog().Fatal("Update index error")
		}
	}

	currencyList := new(CurrencyList)
	this.currency = currencyList.Get()
	if 0 == len(this.currency) {
		currencyList.Update(nil)
		this.currency = currencyList.Get()
		if 0 == len(this.currency) {
			entity.GetLog().Fatal("Update currency error")
		}
	}

	ids := this.getAllId()
	sliHtml := new(FundHtml).Update(ids)
	sliJs := new(FundJs).Update(ids)
	sliRate := new(Rate).Update(ids)

	sliHtml = append(sliHtml, sliJs...)
	sliHtml = append(sliHtml, sliRate...)
	setTask(sliHtml)
}
