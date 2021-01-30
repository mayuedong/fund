package search

import (
	"fund/entity"
	"fund/util"
)

type (
	investNode struct {
		Id      string  `json:"指数代码"`
		Name    string  `json:"指数名称"`
		Price   int     `json:"申购金额"`
		Percent float64 `json:"浮动百分比"`
		Cur     float64 `json:"当前点"`
		Build   float64 `json:"建仓点"`
		Plus    float64 `json:"加仓点"`
		Max     float64 `json:"历史最高点"`
		Min     float64 `json:"历史最低点"`
	}
	invest struct {
		Price    int           `json:"申购金额"`
		Turnover float64       `json:"A股成交额"`
		Indexs   []*investNode `json:"建仓"`
		Pool     []*investNode `json:"等待"`
	}
)

func (this *invest) search(req *request) {
	ptrCurIndex := util.GetCurIndex()
	this.Turnover = ptrCurIndex.GetTurnover(`000001`) / 1e8
	topics := entity.GetConf().GetIndexTopics()
	sli := util.GetHisInfo()
	for _, v := range sli {
		node := new(investNode)
		node.Id = v.Id
		node.Name = topics[v.Id]
		node.Min = v.Asc[v.Min].Index
		node.Max = v.Asc[v.Max].Index
		node.Build = v.Asc[v.Build].Index
		node.Plus = v.Asc[v.Plus].Index
		node.Cur = ptrCurIndex.GetPrice(v.Id)
		node.Percent = ptrCurIndex.GetPercent(v.Id) * 100
		node.Price = int(v.GetPrice(node.Cur))
		if int(entity.GetConf().GetBasePrice()) <= node.Price {
			this.Price += node.Price
			this.Indexs = append(this.Indexs, node)
		} else {
			this.Pool = append(this.Pool, node)
		}
	}
}
