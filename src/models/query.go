package models

import "strings"

type Query struct {
	Total  uint16
	Offset uint8
	Sort   string
}

type Filter struct {
	Author   []uint32
	Editor   []uint32
	Category []uint32
	Taxonomy []uint32
	Status   string
	Total    uint16
	Offset   uint8
	Sort     []string
}

func NewQuery(total uint16, offset uint8, sort ...string) Query {
	query := Query{
		Total:  total,
		Offset: offset,
		Sort:   strings.Join(sort, " "),
	}

	if query.Total < 0 || query.Total > 100 {
		query.Total = 10
	}

	if query.Offset < 0 {
		query.Offset = 0
	}

	if len(query.Sort) == 0 {
		query.Sort = "CreatedAt"
	}

	return query
}
