package consensus

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"
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

func (ba *Baxos) Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool, num_replicas_chan chan int, process_name_chan chan string) {
	replica_path := "/bench/replica"
	ctl_path := "/bench/client"

	num_replicas, err := strconv.ParseInt(ba.options.Option["num_replicas"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_replicas")

	}
	num_clients, err := strconv.ParseInt(ba.options.Option["num_clients"], 10, 64)
	if err != nil {
		panic(err.Error() + " while parsing num_clients")
	}

	round_trip_time, ok := ba.options.Option["round_trip_time"]
	if !ok {
		panic(err.Error() + " while parsing round_trip_time")
	}

	arrival_rate, ok := ba.options.Option["arrival_rate"]
	if !ok {
		panic(err.Error() + " while parsing arrival_rate")
	}

	ports, ok := ba.options.Option["ports"]
	if !ok {
		panic(err.Error() + " while parsing ports")
	}

	var wg sync.WaitGroup
	wg.Add(int(num_replicas + num_clients))
	for i := 0; i < int(num_replicas+num_clients); i++ {
		go func(j int) {
			nodes[j].ExecCmd("pkill replica")
			nodes[j].ExecCmd("pkill client")
			nodes[j].ExecCmd(fmt.Sprintf("rm -r %vbench/logs/", nodes[j].HomeDir))
			nodes[j].ExecCmd(fmt.Sprintf("mkdir -p %vbench/logs/", nodes[j].HomeDir))
			wg.Done()
		}(i)
	}
	wg.Wait()

	fmt.Print("Killed all the replicas and clients\n")
	num_replicas_chan <- int(num_replicas)
	process_name_chan <- "replica" + "," + ports

	for j := 0; j < int(num_replicas); j++ {
		go func(i int) {
			nodes[i].ExecCmd("." + replica_path + " --name " + strconv.Itoa(i+1) + " --roundTripTime " + round_trip_time + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[i].HomeDir))
		}(j)
	}

	time.Sleep(15 * time.Second)

	fmt.Print("Started all the replicas\n")

	nodes[num_replicas].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(50+1) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[num_replicas].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[num_replicas].HomeDir) + " --requestType status --operationType 1 ")

	fmt.Print("Sent initial status to bootstrap\n")

	time.Sleep(20 * time.Second)

	clientOutputs := make([]string, num_clients)
	m := 1
	for j := int(num_replicas); j < int(num_replicas+num_clients); j++ {
		go func(i int, k int) {
			clientOutputs[i-int(num_replicas)] = nodes[i].ExecCmd("." + ctl_path + " --name " + strconv.Itoa(50+k) + " --logFilePath " + fmt.Sprintf("%vbench/logs/", nodes[i].HomeDir) + " --config " + fmt.Sprintf("%vbench/ip_config.yaml", nodes[i].HomeDir) + " --requestType request --arrivalRate  " + arrival_rate + " --testDuration " + strconv.Itoa(duration))
		}(j, m)
		m++
	}

	fmt.Print("Started all the clients\n")

	time.Sleep(6 * time.Second)

	fmt.Print("Bootstrap complete\n")
	bootstrap_complete <- true

	time.Sleep(time.Duration(2*duration) * time.Second)

	fmt.Print("Finished the clients\n")

	var wg1 sync.WaitGroup
	wg1.Add(int(num_replicas + num_clients))
	for j := 0; j < int(num_replicas+num_clients); j++ {
		go func(i int) {
			nodes[i].ExecCmd("pkill replica")
			nodes[i].ExecCmd("pkill client")
			wg1.Done()
		}(j)
	}
	wg1.Wait()

	fmt.Print("Killed all the replicas and clients\n")

	fmt.Printf("Client outputs: %v\n", clientOutputs)

	result <- util.Performance{} //TODO
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
