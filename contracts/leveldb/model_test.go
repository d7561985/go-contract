// +build unit

package leveldb

import (
	"reflect"
	"testing"
)

func TestSimpleQuery_Filter(t *testing.T) {
	type args struct {
		f Filter
	}
	tests := []struct {
		name    string
		sl      SimpleQuery
		args    args
		wantRes SimpleQuery
		wantErr bool
	}{
		{
			name: "filter string",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "BY"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "RU"}}},
			},
			args: args{f: Filter{Key: "country", Value: "BY"}},
			wantRes: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "BY"}}},
			},
			wantErr: false,
		},
		{
			name: "filter int",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"num": 1}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 2}}},
			},
			args: args{f: Filter{Key: "num", Value: "1"}},
			wantRes: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"num": 1}}},
			},
			wantErr: false,
		},
		{
			name: "filter float",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"num": 1.1}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 2.2}}},
			},
			args: args{f: Filter{Key: "num", Value: "1.1"}},
			wantRes: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"num": 1.1}}},
			},
			wantErr: false,
		},
		{
			name: "filter bool",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"b": true}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"b": false}}},
			},
			args: args{f: Filter{Key: "b", Value: "true"}},
			wantRes: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"b": true}}},
			},
			wantErr: false,
		},
		{
			name: "filter bad param",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"num": 1}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 2}}},
			},
			args: args{f: Filter{Key: "num", Value: "XXX"}},
			wantRes: nil,
			wantErr: true,
		},
		{
			name: "filter float bad format",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"num": 1.1}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 2.2}}},
			},
			args: args{f: Filter{Key: "num", Value: "XXX"}},
			wantRes: nil,
			wantErr: true,
		},
		{
			name:    "filter empty",
			sl:      []Query{},
			args:    args{f: Filter{Key: "country", Value: "BY"}},
			wantRes: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.sl.Filter(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Filter() gotRes = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestSimpleQuery_Sort(t *testing.T) {
	type args struct {
		s Sort
	}
	tests := []struct {
		name    string
		sl      SimpleQuery
		args    args
		want    SimpleQuery
		wantErr bool
	}{
		{
			name: "sort string",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": "BB"}}},
			},
			args: args{s: Sort{Field: "country", Asc: true}},
			want: []Query{
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": "BB"}}},
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},

			},
			wantErr: false,
		},
		{
			name: "sort string desc",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": "BB"}}},
			},
			args: args{s: Sort{Field: "country", Asc: false}},
			want: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": "BB"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},

			},
			wantErr: false,
		},
		{
			name: "sort string some field not exists",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "4", Object: SimpleQueue{Context: map[string]interface{}{"num": "5"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": "BB"}}},
			},
			args: args{s: Sort{Field: "country", Asc: true}},
			want: []Query{
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": "BB"}}},
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "4", Object: SimpleQueue{Context: map[string]interface{}{"num": "5"}}},
			},
			wantErr: false,
		},
		{
			name: "sort string desc some field not exists",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "4", Object: SimpleQueue{Context: map[string]interface{}{"num": "5"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": "BB"}}},
			},
			args: args{s: Sort{Field: "country", Asc: false}},
			want: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": "BB"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "4", Object: SimpleQueue{Context: map[string]interface{}{"num": "5"}}},
			},
			wantErr: false,
		},
		{
			name: "sort string different types asc",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": 1}}},
			},
			args: args{s: Sort{Field: "country", Asc: true}},
			want: []Query{
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": 1}}},
			},
			wantErr: false,
		},
		{
			name: "sort string different types desc",
			sl: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": 1}}},
			},
			args: args{s: Sort{Field: "country", Asc: false}},
			want: []Query{
				{Key: "0", Object: SimpleQueue{Context: map[string]interface{}{"country": "ZZ"}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"country": "AA"}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"country": 1}}},
			},
			wantErr: false,
		},
		{
			name: "sort num",
			sl: []Query{
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": 3}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 1}}},
				{Key: "2", Object: SimpleQueue{Context: map[string]interface{}{"num": 2}}},
			},
			args: args{s: Sort{Field: "num", Asc: true}},
			want: []Query{
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 1}}},
				{Key: "2", Object: SimpleQueue{Context: map[string]interface{}{"num": 2}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": 3}}},

			},
			wantErr: false,
		},
		{
			name: "sort num desc",
			sl: []Query{
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": 3}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 1}}},
				{Key: "2", Object: SimpleQueue{Context: map[string]interface{}{"num": 2}}},
			},
			args: args{s: Sort{Field: "num", Asc: false}},
			want: []Query{
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": 3}}},
				{Key: "2", Object: SimpleQueue{Context: map[string]interface{}{"num": 2}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 1}}},

			},
			wantErr: false,
		},
		{
			name: "sort float",
			sl: []Query{
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": 3.1}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 1.2}}},
				{Key: "2", Object: SimpleQueue{Context: map[string]interface{}{"num": 2.3}}},
			},
			args: args{s: Sort{Field: "num", Asc: true}},
			want: []Query{
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 1.2}}},
				{Key: "2", Object: SimpleQueue{Context: map[string]interface{}{"num": 2.3}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": 3.1}}},

			},
			wantErr: false,
		},
		{
			name: "sort float desc",
			sl: []Query{
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": 3.1}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 1.2}}},
				{Key: "2", Object: SimpleQueue{Context: map[string]interface{}{"num": 2.3}}},
			},
			args: args{s: Sort{Field: "num", Asc: false}},
			want: []Query{
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": 3.1}}},
				{Key: "2", Object: SimpleQueue{Context: map[string]interface{}{"num": 2.3}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": 1.2}}},

			},
			wantErr: false,
		},
		{
			name: "sort bool",
			sl: []Query{
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": true}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": false}}},
			},
			args: args{s: Sort{Field: "num", Asc: true}},
			want: []Query{
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": true}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": false}}},
			},
			wantErr: false,
		},
		{
			name: "sort bool desc",
			sl: []Query{
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": true}}},
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": false}}},
			},
			args: args{s: Sort{Field: "num", Asc: false}},
			want: []Query{
				{Key: "1", Object: SimpleQueue{Context: map[string]interface{}{"num": false}}},
				{Key: "3", Object: SimpleQueue{Context: map[string]interface{}{"num": true}}},

			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.sl.Sort(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sort() got = %v, want %v", got, tt.want)
			}
		})
	}
}
