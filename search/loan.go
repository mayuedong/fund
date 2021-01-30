package search

import (
	"fund/entity"
	"math"
	"strconv"
)

type (
	loanMonth struct {
		Num   float64
		Bj    float64
		Lx    float64
		Bx    float64
		SubBj float64
	}
	loanList struct {
		Money float64
		Years float64
		Bx    float64
	}
)

func loanBx(money, month, rate float64) float64 {
	sumrate := 0.0
	for i := 0.0; i < month; i += 1.0 {
		sumrate += math.Pow(1.0+rate, i)
	}
	monthBj := money / sumrate
	monthLx := money * rate
	return monthBj + monthLx
}

func LoanDetail(req *request) (sli []*loanMonth) {
	money, err := strconv.ParseFloat(req.Money, 10)
	if nil != err {
		entity.GetLog().Println(err)
		return
	}
	money *= 1e4

	years, err := strconv.ParseFloat(req.Years, 10)
	if nil != err {
		entity.GetLog().Println(err)
		return
	}
	months := years * 12

	sumrate := 0.0
	rate := entity.GetConf().GetRate() / 1200.0
	for i := 0.0; i < months; i += 1.0 {
		sumrate += math.Pow(1.0+rate, i)
	}

	monthBj := money / sumrate
	monthLx := money * rate
	bx := monthBj + monthLx
	for i := 0.0; i < months; i += 1 {
		p := new(loanMonth)
		p.Num = i + 1
		p.Bj = monthBj * math.Pow(1+rate, i)
		p.Bx = bx
		p.Lx = bx - p.Bj
		if 0 == len(sli) {
			p.SubBj = money - p.Bj
		} else {
			p.SubBj = sli[len(sli)-1].SubBj - p.Bj
		}
		sli = append(sli, p)
	}
	return sli
}

func LoanList(req *request) (sli []*loanList) {
	rate := entity.GetConf().GetRate()
	limit := entity.GetConf().GetLimit()
	for m := 50.0; m <= 100.0; m += 5.0 {
		for y := 10.0; y <= 30.0; y += 5 {
			bx := loanBx(m*1e4, y*12, rate/1200.0)
			if limit-500 < bx && bx < limit+500 {
				list := new(loanList)
				list.Money = m
				list.Years = y
				list.Bx = bx
				sli = append(sli, list)
			}
		}
	}
	return sli
}
