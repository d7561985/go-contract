// +build unit

package leveldb

import (
	"reflect"
	"testing"
)

func TestParseOperation(t *testing.T) {
	type args struct {
		op string
	}
	tests := []struct {
		name    string
		args    args
		want    *operation
		wantErr bool
	}{
		{
			"full",
			args{op: "from=0&to=1558080533-758077000&sort=country&filter=country=BY"},
			&operation{
				Selector: Selector{
					From: "0",
					To:   "1558080533-758077000",
				},
				Filter: Filter{
					Key:   "country",
					Value: "BY",
				},
				Sort: Sort{
					Field: "country",
					Asc:   true,
				},
			},
			false,
		},
		{
			"full-desc",
			args{op: "from=0&to=1558080533-758077000&sort=-country&filter=country=BY"},
			&operation{
				Selector: Selector{
					From: "0",
					To:   "1558080533-758077000",
				},
				Filter: Filter{
					Key:   "country",
					Value: "BY",
				},
				Sort: Sort{
					Field: "country",
					Asc:   false,
				},
			},
			false,
		},
		{
			"filter-only",
			args{op: "filter=country=BY"},
			&operation{
				Filter: Filter{
					Key:   "country",
					Value: "BY",
				},
			},
			false,
		},
		{
			"filter-bad-format",
			args{op: "filter=country"},
			nil,
			true,
		},
		{
			"filter empty",
			args{op: "filter="},
			&operation{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOperation(tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseOperation() got = %v, want %v", got, tt.want)
			}
		})
	}
}
