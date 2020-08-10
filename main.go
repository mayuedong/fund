package main

import (
	"fund/entity"
	"fund/search"
	"fund/util"
	"net/http"
	"os"
	"time"
)

func main() {
	if e := entity.LoadConf(os.Args[1]); nil != e {
		entity.GetLog().Print(e)
	}

	fundMysql := util.GetFundMysql()
	defer fundMysql.Close()

	server := &http.Server{
		Addr:         "0.0.0.0:9527",
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Handler:      new(search.Search),
	}
	server.ListenAndServe()
}
