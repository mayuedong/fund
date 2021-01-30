package search

import (
	"encoding/json"
	"fund/entity"
	"fund/util"
	"net/http"
)

type Search bool

func (this *Search) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := new(request)
	req.parseQuery(r.URL.RawQuery)
	entity.GetLog().Print(r.URL.RawQuery)

	var ret interface{}
	switch req.getSt() {
	case "invest":
		ptr := new(invest)
		ptr.search(req)
		ret = ptr
	case "mix":
		ptr := new(mixPool)
		ptr.search(req)
		ret = ptr
	case "index":
		ptr := new(indexPool)
		ptr.search(req)
		ret = ptr
	case "currency":
		ptr := new(currencyPool)
		ptr.search(req)
		ret = ptr
	case "test":
		ptr := new(util.HisInfo)
		ret = ptr.Test()
	case "loanList":
		ret = LoanList(req)
	case "loanDetail":
		ret = LoanDetail(req)
	}
	sliByte, err := json.Marshal(ret)
	if nil != err {
		entity.GetLog().Print(err)
	}
	w.Write([]byte(sliByte))
}
