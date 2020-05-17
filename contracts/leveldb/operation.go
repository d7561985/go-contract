package leveldb

import (
	"fmt"
	"net/url"
	"strings"
)

type Filter struct {
	Key   string
	Value string
}

type Selector struct {
	From string
	To   string
}

type Sort struct {
	Field string
	Asc   bool
}

type operation struct {
	Selector Selector
	Filter   Filter
	Sort     Sort
}

const (
	QuerySelectorFrom = "from"
	QuerySelectorTo   = "to"
	QueryFilter       = "filter"
	QuerySort         = "sort"
)

// ParseOperation as url query
// Selector: from, to
// Filter: Filter
// Sort: Sort, argument prefix support - DESC
func ParseOperation(op string) (*operation, error) {
	q, err := url.ParseQuery(op)
	if err != nil {
		return nil, fmt.Errorf("parse operation error: %w", err)
	}

	res := &operation{}

	for s, vals := range q {
		if vals[0] == "" {
			continue
		}

		switch s {
		case QuerySelectorFrom:
			res.Selector.From = vals[0]
		case QuerySelectorTo:
			res.Selector.To = vals[0]
		case QueryFilter:
			f := strings.Split(vals[0], "=")
			if len(f) != 2 {
				return nil, fmt.Errorf(`wrong Filter format. support only context field equasion. example: "country=RU"`)
			}

			res.Filter.Key = f[0]
			res.Filter.Value = f[1]
		case QuerySort:
			if vals[0][0] != '-'{
				res.Sort.Asc = true
			}

			res.Sort.Field = strings.TrimPrefix(vals[0], "-")
		}
	}

	return res, nil
}
