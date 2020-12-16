package util

import (
	"fund/entity"
)

type Update struct {
	mix      []*MixList
	index    []*IndexList
	currency []*CurrencyList
	turnover []*Turnover
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
	this.turnover = history.Get()
	if 0 == len(this.turnover) {
		history.Update(nil)
		this.turnover = history.Get()
		if 0 == len(this.turnover) {
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
	this.delOverdue(fundHtmlTable, ids)
	new(FundHtml).Update(ids)

	this.delOverdue(fundJsTable, ids)
	new(FundJs).Update(ids)

	this.delOverdue(rateTable, ids)
	new(Rate).Update(ids)
}

func (this *Update) delOverdue(table string, ids []string) {
	db := GetSqlite(table)
	defer db.CLOSE()
	keys := db.KEYS()
	for _, key := range keys {
		i := 0
		for ; i < len(ids); i++ {
			if key == ids[i] {
				break
			}
		}
		if i == len(ids) {
			db.DELETE(key)
		}
	}
}
