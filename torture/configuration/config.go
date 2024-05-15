package configuration

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

/*
	config.go defines the structs and methods to pass the configuration file, that contains the IP:ports of each torture torture
*/

type ReplicaInstance struct {
	Name          string
	IP            string
	PORT          string
	REPLICA_PORTS []string
}

type Controller struct {
	Name string
	IP   string
	PORT string
}

// InstanceConfig describes the set of replicas
type InstanceConfig struct {
	Peers      []ReplicaInstance
	Controller Controller
}

// NewInstanceConfig loads an instance configuration from given file
func NewInstanceConfig(fname string, name int64) (*InstanceConfig, error) {
	cfg := InstanceConfig{
		Peers: []ReplicaInstance{},
	}

	file, err := os.Open(fname)
	if err != nil {
		panic(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err.Error())
		}
	}(file)

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	for i := 0; i < len(lines)-1; i++ {
		line := lines[i]
		parts := strings.Split(line, " ")

		// Create a new ReplicaInstance
		peer := ReplicaInstance{
			Name:          parts[0],
			IP:            parts[1],
			PORT:          parts[2],
			REPLICA_PORTS: parts[3:],
		}

		// Append the new ReplicaInstance to the configuration
		cfg.Peers = append(cfg.Peers, peer)
	}
	c_line := strings.Split(lines[len(lines)-1], " ")
	c := Controller{
		Name: c_line[0],
		IP:   c_line[1],
		PORT: c_line[2],
	}

	cfg.Controller = c

	// set the self ip to 0.0.0.0
	cfg = configureSelfIP(cfg, name)
	return &cfg, nil
}

/*
	Replace the IP of my self to 0.0.0.0
*/

func configureSelfIP(cfg InstanceConfig, name int64) InstanceConfig {
	for i := 0; i < len(cfg.Peers); i++ {
		if cfg.Peers[i].Name == strconv.Itoa(int(name)) {
			cfg.Peers[i].IP = "0.0.0.0"
			return cfg
		}
	}
	return cfg
}
