package util

import (
	"fund/entity"
	"strconv"
	"strings"
)

type IndexList struct {
	Id    string
	Name  string
	Year1 float64
	Year2 float64
	Year3 float64
}

func (r *IndexList) Get() (sli []*IndexList) {
	db := GetSqlite(indexListTable)
	defer db.CLOSE()
	keys := db.KEYS()
	for _, key := range keys {
		ptr := new(IndexList)
		db.GET(key, ptr)
		sli = append(sli, ptr)
	}
	return sli
}
func (r *IndexList) GetOne(id string) *IndexList {
	db := GetSqlite(indexListTable)
	defer db.CLOSE()
	ptr := new(IndexList)
	db.GET(id, ptr)
	return ptr
}

func (r *IndexList) Update(ids []string) {
	Download(new(IndexList))
}

func (this *IndexList) getUptime() string {
	return ""
}

func (this *IndexList) getWait() int {
	return -30
}

func (this *IndexList) getUrl() string {
	return entity.GetConf().GetIndexList()
}
func (this *IndexList) parse(b []byte) error {
	sliFund := parse(b, `["`, `"]`, `","`)
	if 0 == len(sliFund) {
		entity.GetLog().Print("pase IndexList ", string(b))
	}
	db := GetSqlite(indexListTable)
	defer db.CLOSE()
	for _, v := range sliFund {
		if sli := strings.Split(v, "|"); 13 > len(sli) {
			entity.GetLog().Print("parse Mix err : ", v)
			continue
		} else {
			ptr := new(IndexList)
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
