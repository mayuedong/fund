package search

import (
	"fund/entity"
	"fund/util"
	"time"
)

type update struct {
	mixFund   []*util.FundList
	indexFund []*util.FundList
}

func (this *update) filterMixAchievement() {
	sli := this.mixFund
	l := int(len(sli) / 4)
	var good1, good2, good3, good []*util.FundList
	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].GetYear1() < sli[j].GetYear1() {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	good1 = append(good1, sli[:l]...)

	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].GetYear2() < sli[j].GetYear2() {
				sli[i], sli[j] = sli[j], sli[i]
			}
		}
	}
	good2 = append(good2, sli[:l]...)

	for i := 0; i < len(sli); i++ {
		for j := i + 1; j < len(sli); j++ {
			if sli[i].GetYear3() < sli[j].GetYear3() {
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
	this.mixFund = good
}

func (this *update) filterScale() {
	var tmp []*util.FundList
	ptrInfo := util.GetFundInfo()
	for _, ptr := range this.mixFund {
		if 35.0 < ptrInfo.GetScale(ptr.GetId()) {
			tmp = append(tmp, ptr)
		}
	}
	this.mixFund = tmp

	var aux []*util.FundList
	for _, ptr := range this.indexFund {
		if 3.0 < ptrInfo.GetScale(ptr.GetId()) {
			aux = append(aux, ptr)
		}
	}
	this.indexFund = aux
}

var g_updateIntervalSecond int64

func (this *update) search() {
	now := time.Now()
	if 0 == g_updateIntervalSecond {
		g_updateIntervalSecond = now.Unix()
	} else {
		curUnix := now.Unix()
		if curUnix-g_updateIntervalSecond < 3600 {
			return
		}
	}

	ptrHisIndex := util.GetHisIndex()
	ptrHisIndex.Update()

	ptrFund := util.GetFund()
	ptrFund.Update()

	this.indexFund = ptrFund.GetIndex()
	this.mixFund = ptrFund.GetMix()
	this.filterMixAchievement()

	var infoIds []string
	for _, v := range this.indexFund {
		infoIds = append(infoIds, v.GetId())
	}
	for _, v := range this.mixFund {
		infoIds = append(infoIds, v.GetId())
	}
	ptrFundInfo := util.GetFundInfo()
	ptrFundInfo.Update(infoIds)
	this.filterScale()

	var mixInfoIds []string
	for _, v := range this.mixFund {
		mixInfoIds = append(mixInfoIds, v.GetId())
	}
	for _, v := range this.indexFund {
		mixInfoIds = append(mixInfoIds, v.GetId())
	}
	ptrMixInfo := util.GetMixInfo()
	ptrMixInfo.Update(mixInfoIds)

	ptrCostRate := util.GetCostRate()
	index := new(indexPool)
	index.search(nil)
	indexIds := index.getIds()
	ptrCostRate.Update(indexIds)
	entity.GetLog().Print("successful")
	return
}
