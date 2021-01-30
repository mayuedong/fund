package search

import (
	"net/url"
	"strings"
)

type request struct {
	St    string
	Money string
	Years string
}

func (this *request) getSt() string {
	return this.St
}

func (this *request) parseQuery(inQuery string) {
	mapQuery := make(map[string]string)
	for _, kv := range strings.Split(inQuery, "&") {
		pos := strings.Index(kv, "=")
		if 0 < pos && pos+1 < len(kv) {
			mapQuery[kv[:pos]], _ = url.QueryUnescape(kv[pos+1:])
		}
	}
	this.St = mapQuery["st"]
	this.Money = mapQuery["money"]
	this.Years = mapQuery["years"]
}
