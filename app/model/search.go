package model

import (
	"fmt"
	"strings"
)

type SearchData struct {
	Search string `json:"-" schema:"search"`
}

func (d *SearchData) CheckSearchData() bool {
	return len(d.Search) > 0
}

func (d *SearchData) SearchLikes(fields []string) string {
	res := ""
	if d.CheckSearchData() {
		for _, l := range fields {
			search := strings.ToLower(d.Search)
			res = fmt.Sprintf("%s OR lower(%s::varchar(255)) LIKE '%%%s%%'", res, l, strings.ReplaceAll(search, "'", "%"))
		}
	}

	return strings.TrimPrefix(res, " OR ")
}
