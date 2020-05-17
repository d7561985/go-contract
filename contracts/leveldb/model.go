package leveldb

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Context map[string]interface{}

// SimpleQueue uses time as key identification for queue
// support only valid JSON extra-context
type SimpleQueue struct {
	// this property should strictly follow UTF8 Sort rule for identify assets
	Time time.Time `json:"created_at"`

	Context Context `json:"context"`
}

// NewSimpleQueue create empty queue element
func NewSimpleQueue() SimpleQueue {
	return SimpleQueue{Time: time.Now(), Context: make(Context)}
}

func (s *SimpleQueue) BLOB() ([]byte, error) {
	return json.Marshal(s)
}

type Query struct {
	Key    string      `json:"key"`
	Object SimpleQueue `json:"object"`
}

type SimpleQuery []Query

// Key generate unique queue value.
// For sake of decreasing collision we use nanosecond postfix
// But collision possible, for now it's low chance but for handling that we should make a lot more stuff
func TimedKey(time time.Time) string {
	return fmt.Sprintf("%d-%d", time.Unix(), time.Nanosecond())
}

// small helper
func mustParse(value string) time.Time {
	v, _ := time.Parse(time.RFC3339Nano, value)
	return v
}

// Filter return filtered data.
// Currently supported types: float64, int, string, bool
// Just in case: marshaling all numbers transform into float64
func (sl SimpleQuery) Filter(f Filter) (res SimpleQuery, err error) {
	for i, queue := range sl {
		v, ok := queue.Object.Context[f.Key]
		if !ok {
			continue
		}

		add := false

		switch exp := v.(type) {
		case string:
			add = exp == f.Value

		case bool:
			add = (exp && strings.ToUpper(f.Value) == "TRUE") ||
				(!exp && strings.ToUpper(f.Value) == "FALSE")

		case int:
			p, err := strconv.ParseInt(f.Value, 0, 64)
			if err != nil {
				return nil, fmt.Errorf("can't parse %q to float64 error: %w", f.Value, err)
			}

			add = int(p) == exp
		case float64:
			p, err := strconv.ParseFloat(f.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("can't parse %q to float64 error: %w", f.Value, err)
			}

			add = p == exp
		default:
			log.Printf("Filter unsuported type %T for key %q val %q", v, f.Key, f.Value)
			continue
		}

		if add {
			res = append(res, sl[i])
		}

	}

	return res, nil
}

// Sort query array
func (sl SimpleQuery) Sort(s Sort) (SimpleQuery, error) {
	if s.Field == "" {
		return sl, nil
	}

	sort.Slice(sl, func(i, j int) bool {
		vi, ok := sl[i].Object.Context[s.Field]
		if !ok {
			return false
		}

		vj, ok := sl[j].Object.Context[s.Field]
		if !ok {
			return true
		}

		switch I := vi.(type) {
		case string:
			J, ok := vj.(string)
			if !ok {
				log.Printf("sort key %q has context with different types %T and %T", s.Field, vi, vj)
				return false
			}

			// I better prefer XOR operation
			switch {
			case s.Asc:
				return I < J
			default:
				return I > J
			}
		case int:
			J, ok := vj.(int)
			if !ok {
				log.Printf("sort key %q has context with different types %T and %T", s.Field, vi, vj)
				return false
			}

			// I better prefer XOR operation
			switch {
			case s.Asc:
				return I < J
			default:
				return I > J
			}
		case float64:
			J, ok := vj.(float64)
			if !ok {
				log.Printf("sort key %q has context with different types %T and %T", s.Field, vi, vj)
				return false
			}

			// I better prefer XOR operation
			switch {
			case s.Asc:
				return I < J
			default:
				return I > J
			}
		case bool:
			J, ok := vj.(bool)
			if !ok {
				log.Printf("sort key %q has context with different types %T and %T", s.Field, vi, vj)
				return false
			}

			// I better prefer XOR operation
			switch {
			case s.Asc:
				return !J && I
			default:
				return J && !I
			}

		default:
			return false
		}
	})

	return sl, nil
}
