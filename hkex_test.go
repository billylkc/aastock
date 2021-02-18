package aastock

import (
	"reflect"
	"testing"
)

func TestGetCompanyName(t *testing.T) {
	type args struct {
		c int
	}
	tests := []struct {
		name    string
		args    args
		want    company
		wantErr bool
	}{
		{
			name: "Simple test",
			args: args{
				c: 5, // 00005 as hsbc
			},
			want: company{
				Code: "00005",
				Name: "匯豐控股",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCompanyName(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCompanyName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCompanyName() = %v, want %v", got, tt.want)
			}
		})
	}
}
