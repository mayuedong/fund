package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fund/entity"
	"fund/util"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Search bool

func (this *Search) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		b, _ := ioutil.ReadAll(r.Body)
		var m map[string][]byte
		err := json.Unmarshal(b, &m)
		if nil != err {
			entity.GetLog().Fatal(err)
		}
		reader := bytes.NewReader(m["image"])
		img, err := png.Decode(reader)
		if nil != err {
			entity.GetLog().Fatal(err)
		}

		fName := fmt.Sprintf("myd_%s.png", time.Now().String())
		out, err := os.Create(fName)
		if err != nil {
			entity.GetLog().Fatal(err)
		}
		defer out.Close()

		err = png.Encode(out, img)
		if nil != err {
			entity.GetLog().Fatal(err)
		}
		r.Body.Close()
	}

	req := new(request)
	req.parseQuery(r.URL.RawQuery)

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
	default:
		ptr := new(util.Turnover)
		ret = ptr.Test()
	}
	sliByte, err := json.Marshal(ret)
	if nil != err {
		entity.GetLog().Print(err)
	}
	w.Write([]byte(sliByte))
}
