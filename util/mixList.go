package util

import (
	"fund/entity"
	"strconv"
	"strings"
)

type MixList struct {
	Id    string
	Name  string
	Year1 float64
	Year2 float64
	Year3 float64
}

func (r *MixList) Get() (sli []*MixList) {
	db := GetSqlite(mixListTable)
	defer db.CLOSE()
	keys := db.KEYS()
	for _, key := range keys {
		ptr := new(MixList)
		db.GET(key, ptr)
		sli = append(sli, ptr)
	}
	return sli
}

func (r *MixList) Update(ids []string) {
	Download(new(MixList))
}

func (this *MixList) getUptime() string {
	return ""
}

func (this *MixList) getWait() int {
	return -30
}

func (this *MixList) getUrl() string {
	return entity.GetConf().GetMixList()
}

func (this *MixList) parse(b []byte) error {
	sliFund := parse(b, `["`, `"]`, `","`)
	if 0 == len(sliFund) {
		entity.GetLog().Print("pase MixList ", string(b))
	}

	db := GetSqlite(mixListTable)
	defer db.CLOSE()
	for _, v := range sliFund {
		if sli := strings.Split(v, "|"); 13 > len(sli) {
			entity.GetLog().Print("parse MixList err : ", v)
			continue
		} else {
			ptr := new(MixList)
			ptr.Id = sli[0]
			ptr.Name = sli[1]
			ptr.Year1, _ = strconv.ParseFloat(sli[10], 64)
			ptr.Year2, _ = strconv.ParseFloat(sli[11], 64)
			ptr.Year3, _ = strconv.ParseFloat(sli[12], 64)
			db.SET(ptr.Id, ptr)
		}
	}
	return nil
}
