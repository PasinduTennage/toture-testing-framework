package cmd

import (
	"testing"
)

// Test_getProcessIds tests the GetProcessIds function
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
			name: "simple test",
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
			t.Logf("got: %v\n", got)
		})
	}
}

// Test_getProcessID tests the getProcessID function
// before running the test, run the dummy server with --name 1
func Test_getProcessID(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "simple test",
			args: args{
				port: 10001,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetProcessID(tt.args.port)
			if got == -1 {
				t.Errorf("getProcessID() = %v", got)
			} else {
				t.Logf("got: %v\n", got)
			}
		})
	}
}
