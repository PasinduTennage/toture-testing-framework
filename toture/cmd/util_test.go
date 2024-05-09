package cmd

import (
	"reflect"
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
		{
			name: "simple test",
			args: args{
				port: 11001,
			},
		},
		{
			name: "simple test",
			args: args{
				port: 12001,
			},
		},
		{
			name: "simple test",
			args: args{
				port: 13001,
			},
		},
		{
			name: "simple test",
			args: args{
				port: 14001,
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

// Test_runDummyThreads tests the RunDummyThreads function

func TestNewConfig(t *testing.T) {
	type args struct {
		fname string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]int
		wantErr bool
	}{
		{
			name: "simple test",
			args: args{
				fname: "/home/tennage/Documents/toture-testing-consensus/toture/configuration/dummy_config.cfg",
			},
			want: [][]int{
				{10000, 10001, 10002, 10003},
				{11000, 11001, 11002, 11003},
				{12000, 12001, 12002, 12003},
				{13000, 13001, 13002, 13003},
				{14000, 14001, 14002, 14003}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.fname)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInstanceConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInstanceConfig() got = %v, want %v", got, tt.want)
			} else {
				t.Logf("got: %v\n", got)
			}
		})
	}
}
