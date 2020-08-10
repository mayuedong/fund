package search

import (
	"fund/entity"
	"fund/util"
)

type (
	investNode struct {
		Id    string  `json:"指数代码"`
		Price int     `json:"申购金额"`
		Cur   float64 `json:"当前点"`
		Build float64 `json:"建仓点"`
		Plus  float64 `json:"加仓点"`
		Max   float64 `json:"历史最高点"`
		Min   float64 `json:"历史最低点"`
		Rate  float64 `json:"增长率"`
	}
	invest struct {
		Price    int           `json:"申购金额"`
		Turnover float64       `json:"A股成交额"`
		Indexs   []*investNode `json:"指数"`
	}
)

func (this *invest) search(req *request) {
	ptrCurIndex := util.GetCurIndex()
	this.Turnover = ptrCurIndex.GetTurnover(`000001`) / 1e8

	ptrIndexInfo := util.GetIndexInfo()
	codes := entity.GetConf().GetIndexData()
	for _, k := range codes {
		k = k[1:]
		node := new(investNode)
		node.Id = k
		node.Cur = ptrCurIndex.GetPrice(k)
		node.Build = ptrIndexInfo.GetMedium(k)
		node.Plus = ptrIndexInfo.GetSeMe(k)
		node.Price = ptrIndexInfo.GetRate(k, node.Cur)
		node.Max = ptrIndexInfo.GetHigh(k)
		node.Min = ptrIndexInfo.GetLow(k)
		node.Rate = 1.0 - node.Min/node.Max
		this.Price += node.Price
		this.Indexs = append(this.Indexs, node)
	}
}
