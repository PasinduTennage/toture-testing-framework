package consensus

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type Baxos struct {
	logger  *util.Logger
	options protocols.ConsensusOptions
}

func NewBaxos(logger *util.Logger) *Baxos {
	return &Baxos{
		logger: logger,
	}
}

func (ba *Baxos) CopyConsensus(nodes []*common.Node) error {
	num_replicas, ok := ba.options.Option["num_replicas"]
	if !ok {
		panic("num_replicas not found in options")
	}
	num_clients, ok := ba.options.Option["num_clients"]
	if !ok {
		panic("num_clients not found in options")
	}

	num_replicas_int, err := strconv.ParseInt(num_replicas, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	num_clients_int, err := strconv.ParseInt(num_clients, 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_clients")
	}

	if num_replicas_int+num_clients_int > int64(len(nodes)) {
		panic("Not enough nodes to deploy baxos")
	}

	config_inputs := []string{"protocols/baxos/assets/config-generate.py", num_replicas, num_clients}
	for i := int64(0); i < num_clients_int+num_replicas_int; i++ {
		config_inputs = append(config_inputs, nodes[i].Ip)
	}

	fmt.Printf("Running pythong command: %v\n", config_inputs)

	sshCmd := exec.Command("python3", config_inputs...)
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		ba.logger.Debug("Error while running config-generate.py "+err.Error()+" "+string(output), 3)
	}
	if len(output) > 0 {
		// write the output to protocols/baxos/assets/ip_config.yaml
		err = ioutil.WriteFile("protocols/baxos/assets/ip_config.yaml", output, 0644)
		if err != nil {
			panic("Error while writing to ip_config.yaml " + err.Error())
		} else {
			fmt.Printf("ip_config.yaml written successfully with content:\n %v\n", string(output))
		}
	}

	// copy the replica binary, client binary and configuration file to the nodes
	for i := int64(0); i < num_clients_int+num_replicas_int; i++ {
		nodes[i].Put_Load("protocols/baxos/assets/ip_config.yaml", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
		nodes[i].Put_Load("protocols/baxos/assets/replica", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
		nodes[i].Put_Load("protocols/baxos/assets/client", fmt.Sprintf("%vbench/", nodes[i].HomeDir))
	}
	fmt.Print("Copied the baxos binaries to all the nodes\n")

	return nil
}

func (ba *Baxos) Bootstrap(nodes []*common.Node) util.Performance {
	return util.Performance{}
}

func (ba *Baxos) ExtractOptions(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic("Error unmarshalling YAML " + err.Error())
	}

	options := protocols.ConsensusOptions{Option: make(map[string]string)}
	for key, value := range config {
		options.Option[key] = fmt.Sprintf("%v", value)
	}

	fmt.Printf("Baxos options:\n %v\n", options.Option)

	ba.options = options
}
