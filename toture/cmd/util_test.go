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
			name: "sample",
			args: args{
				name: "systemd",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetProcessIds(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getProcessIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("got: %v\n", got)
		})
	}
}
