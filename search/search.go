package search

import (
	"encoding/json"
	"fund/entity"
	"fund/util"
	"io/ioutil"
	"net/http"
)

type Search bool

func (this *Search) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	strReq := ""
	switch r.Method {
	case http.MethodGet:
		strReq = r.URL.RawQuery
	case http.MethodPost:
		byteReq, _ := ioutil.ReadAll(r.Body)
		strReq = string(byteReq)
		r.Body.Close()
	}
	req := new(request)
	req.parseQuery(strReq)
	entity.GetLog().Print(strReq)

	var ret interface{}
	switch req.getSt() {
	case "invest":
		ptr := new(invest)
		ptr.search(req)
		ret = ptr
	case "mixPool":
		ptr := new(mixPool)
		ptr.search(req)
		ret = ptr
	case "indexPool":
		ptr := new(indexPool)
		ptr.search(req)
		ret = ptr
	case "update":
		ptr := new(update)
		ptr.search()
		ret = ptr
	default:
		ptr := util.GetIndexInfo()
		ret = ptr.Test()
	}
	sliByte, err := json.Marshal(ret)
	if nil != err {
		entity.GetLog().Print(err)
	}
	w.Write([]byte(sliByte))
}
