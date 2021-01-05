package aastock

import (
	"reflect"
	"testing"
)

func Test_getLastValue(t *testing.T) {
	type args struct {
		values []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"Normal Case", args{values: []float64{1.23, 2.23, 6.23, 4.23, 6.23, 7.23, 0}}, 7.23},
		{"Normal Case Two", args{values: []float64{1.23, 2.23, 3.23, 0, 0, 0, 4.23, 6.23, 7.23, 0}}, 7.23},
		{"Test for Zeros", args{values: []float64{0, 0, 0, 0}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLastValue(tt.args.values); got != tt.want {
				t.Errorf("getLastValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetList(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentPrice(t *testing.T) {
	type args struct {
		c int
	}
	tests := []struct {
		name    string
		args    args
		want    StockPrice
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCurrentPrice(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCurrentPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
