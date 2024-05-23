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
						Name: "1",
						IP:   "0.0.0.0",
						PORT: "9000",
					},
					{
						Name: "2",
						IP:   "0.0.0.0",
						PORT: "9001",
					},
					{
						Name: "3",
						IP:   "0.0.0.0",
						PORT: "9002",
					},
					{
						Name: "4",
						IP:   "0.0.0.0",
						PORT: "9003",
					},
					{
						Name: "5",
						IP:   "0.0.0.0",
						PORT: "9004",
					},
				},
				Controller: Controller{
					Name: "11",
					IP:   "0.0.0.0",
					PORT: "9999",
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

func TestNewConsensusConfig(t *testing.T) {
	type args struct {
		fname string
	}
	tests := []struct {
		name    string
		args    args
		want    *ConsensusConfig
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				fname: "/home/pasindu/Documents/toture-testing-consensus/torture/configuration/consensus_config/1.cfg",
			},
			want: &ConsensusConfig{
				Options: map[string]string{
					"name":       "1",
					"ip":         "0.0.0.0",
					"ports":      "10000 10001 10002 10003",
					"process_id": "NA",
					"db_port":    "NA",
					"db_file":    "NA",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConsensusConfig(tt.args.fname)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConsensusConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !equalMaps(got.Options, tt.want.Options) {
				t.Errorf("NewConsensusConfig() got = %v, want %v", got, tt.want)
			} else {
				t.Logf("NewConsensusConfig() got = %v", got)
			}
		})
	}
}

// check if 2 maps contain the same key-value pairs
func equalMaps(options1 map[string]string, options2 map[string]string) bool {
	if len(options1) != len(options2) {
		return false
	}

	// Check each key-value pair in the first map
	for key, value1 := range options1 {
		value2, ok := options2[key]
		if !ok || value1 != value2 {
			return false
		}
	}

	// If all checks pass, the maps are equal
	return true
}
