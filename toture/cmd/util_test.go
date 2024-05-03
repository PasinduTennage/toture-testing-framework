package cmd

import (
	"fmt"
	"testing"
)

func Test_getProcessIds(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "chrome",
			args: args{
				name: "chrome",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getProcessIds(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getProcessIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("got: %v\n", got)
		})
	}
}

func Test_getPorts(t *testing.T) {
	type args struct {
		pIds []int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "chrome",
			args: args{
				pIds: []int{2408, 1631, 23079},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPorts(tt.args.pIds)
			fmt.Printf("got: %v\n", got)
		})
	}
}
