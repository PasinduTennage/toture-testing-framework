package configuration

import (
	"reflect"
	"testing"
)

func TestNewInstanceConfig(t *testing.T) {
	type args struct {
		fname string
		name  int64
	}
	tests := []struct {
		name    string
		args    args
		want    *InstanceConfig
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				fname: "/home/pasindu/Documents/toture-testing-consensus/torture/configuration/local-config.cfg",
				name:  1,
			},
			want: &InstanceConfig{
				Peers: []ReplicaInstance{
					{
						Name:          "1",
						IP:            "0.0.0.0",
						PORT:          "9000",
						REPLICA_PORTS: []string{"10000", "10001", "10002", "10003"},
					},
					{
						Name:          "2",
						IP:            "0.0.0.0",
						PORT:          "9001",
						REPLICA_PORTS: []string{"11000", "11001", "11002", "11003"},
					},
					{
						Name:          "3",
						IP:            "0.0.0.0",
						PORT:          "9002",
						REPLICA_PORTS: []string{"12000", "12001", "12002", "12003"},
					},
					{
						Name:          "4",
						IP:            "0.0.0.0",
						PORT:          "9003",
						REPLICA_PORTS: []string{"13000", "13001", "13002", "13003"},
					},
					{
						Name:          "5",
						IP:            "0.0.0.0",
						PORT:          "9004",
						REPLICA_PORTS: []string{"14000", "14001", "14002", "14003"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewInstanceConfig(tt.args.fname, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInstanceConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInstanceConfig() got = %v, want %v", got, tt.want)
			} else {
				t.Logf("NewInstanceConfig() got = %v", got)
			}
		})
	}
}
