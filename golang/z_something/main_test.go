package main

import (
	"reflect"
	"testing"
)

func TestF2(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := F2(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("F2() = %v, want %v", got, tt.want)
			}
		})
	}
}
