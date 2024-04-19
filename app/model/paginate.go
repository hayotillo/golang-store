package model

import "math"

type PaginateData struct {
	PaginatePer  int `json:"per" schema:"per"`
	PaginatePage int `json:"page" schema:"page"`
}

func (d *PaginateData) Per() int {
	if d.PaginatePer == 0 {
		d.PaginatePer = 10
	}
	return d.PaginatePer
}

func (d *PaginateData) Page() int {
	if d.PaginatePage == 0 {
		d.PaginatePage = 1
	}
	return d.PaginatePage
}

func (d *PaginateData) Offset() int {
	return (d.Page() - 1) * d.Per()
}

func (d *PaginateData) Paginate(count int) Paginate {
	res := Paginate{}
	pageCount := int(math.Ceil(float64(count) / float64(d.Per())))

	res.TotalItems = count
	res.Page = d.Page()
	res.TotalPages = pageCount

	return res
}

type Paginate struct {
	TotalItems int `json:"total_items"`
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
}
