package aastock

import (
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
