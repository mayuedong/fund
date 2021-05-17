package search

type (
	fund interface {
		getName() string
		getYear1() float64
		getYear2() float64
		getYear3() float64
		getPrice() float64
		getScale() float64
		getManager() string
		getFoundDate() int
		getFundSize() string
		getTrackAims() string
		getManagerDate() int
	}
	base struct {
		Id          string  `json:"基金代码"`
		Name        string  `json:"基金名称"`
		FundSize    string  `json:"管理规模"`
		Manager     string  `json:"基金经理"`
		Scale       float64 `json:"基金规模"`
		HoldPe      float64 `json:"散户持有百分比"`
		Rate        float64 `json:"运作费用"`
		FoundDate   int     `json:"基金成立天数"`
		ManagerDate int     `json:"基金经理管理天数"`
	}
	currency struct {
		base
		Year1 float64 `json:"近一年收益"`
		Price float64 `json:"起购金额"`
	}
	mix struct {
		base
		UpTime string            `json:"更新时间"`
		Year1  float64           `json:"近一年收益"`
		Year2  float64           `json:"近二年收益"`
		Year3  float64           `json:"近三年收益"`
		Stock  map[string]string `json:"重仓股"`
	}
	index struct {
		base
		TrackAims string            `json:"跟踪标的"`
		TrackErr  float64           `json:"跟踪误差"`
		Year1     float64           `json:"近一年收益"`
		Year2     float64           `json:"近二年收益"`
		Year3     float64           `json:"近三年收益"`
		Stock     map[string]string `json:"重仓股"`
	}
)

func (r *base) getName() string     { return r.Name }
func (r *base) getScale() float64   { return r.Scale }
func (r *base) getManager() string  { return r.Manager }
func (r *base) getManagerDate() int { return r.ManagerDate }
func (r *base) getFoundDate() int   { return r.FoundDate }
func (r *base) getFundSize() string { return r.FundSize }

func (r *mix) getPrice() float64    { return 0.0 }
func (r *mix) getYear1() float64    { return r.Year1 }
func (r *mix) getYear2() float64    { return r.Year2 }
func (r *mix) getYear3() float64    { return r.Year3 }
func (r *mix) getTrackAims() string { return "" }

func (r *index) getPrice() float64    { return 0.0 }
func (r *index) getYear1() float64    { return r.Year1 }
func (r *index) getYear2() float64    { return r.Year2 }
func (r *index) getYear3() float64    { return r.Year3 }
func (r *index) getTrackAims() string { return r.TrackAims }

func (r *currency) getPrice() float64    { return r.Price }
func (r *currency) getYear1() float64    { return r.Year1 }
func (r *currency) getYear2() float64    { return 0.0 }
func (r *currency) getYear3() float64    { return 0.0 }
func (r *currency) getTrackAims() string { return "" }

func filter(sli []fund, fn func(fund) bool) (ret []fund) {
	for _, it := range sli {
		if fn(it) {
			ret = append(ret, it)
		}
	}
	return ret
}

func sort(sli []fund, fn func(max, min fund) bool) (ret []fund) {
	for i := 0; i < len(sli); i++ {
		j := 0
		for ; j < len(ret); j++ {
			if fn(sli[i], ret[j]) {
				ret = append(ret[:j+1], ret[j:]...)
				ret[j] = sli[i]
				break
			}
		}
		if j == len(ret) {
			ret = append(ret, sli[i])
		}
	}
	return ret
}
