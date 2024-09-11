package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/exec"
	"sync"
	"toture-test/util"
)

type NodeStat struct {
	cpu_usage   float64
	mem_usage   float64
	network_in  float64
	network_out float64
}

type Node struct {
	Id             int    `yaml:"Id"`
	Ip             string `yaml:"Ip"`
	Username       string `yaml:"Username"`
	HomeDir        string `yaml:"HomeDir"`
	PrivateKeyPath string `yaml:"privateKeyPath"`
	stat           NodeStat
	statMutex      *sync.Mutex
	Logger         *util.Logger
}

func (n *Node) InitNode(logger *util.Logger) {
	n.stat = NodeStat{
		cpu_usage:   0.0,
		mem_usage:   0.0,
		network_in:  0.0,
		network_out: 0.0,
	}
	n.statMutex = &sync.Mutex{}
	n.Logger = logger
}

// Execute a command on the node

func (n *Node) ExecCmd(cmd string) error {
	sshCmd := exec.Command("ssh", "-i", n.PrivateKeyPath, fmt.Sprintf("%s@%s", n.Username, n.Ip), cmd)
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("failed to execute %v via SSH, err:%v, output:%vs for node:%v", fmt.Sprintf("%v\n", sshCmd), err, string(output), n.Id))
	} else {
		n.Logger.Debug(fmt.Sprintf("Executed %v via SSH, output: %v for node: %v", fmt.Sprintf("%v\n", sshCmd), string(output), n.Id), 0)
	}
	return nil
}

// download the file from the remote location to the local location

func (n *Node) Get_Load(remote_location string, local_location string) error {

	scpCmd := exec.Command("scp", "-i", n.PrivateKeyPath, fmt.Sprintf("%s@%s:%s", n.Username, n.Ip, remote_location), local_location)
	output, err := scpCmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("failed to download file via SCP %v, error:%v, output:%v for node:%v", fmt.Sprintf("%v\n", scpCmd), err, string(output), n.Id))
	} else {
		n.Logger.Debug(fmt.Sprintf("Download using %v successful for node %v", fmt.Sprintf("%v\n", scpCmd), n.Id), 0)
	}
	return nil
}

// upload the file from the local location to the remote location

func (n *Node) Put_Load(local_location string, remote_location string) error {
	scpCmd := exec.Command("scp", "-i", n.PrivateKeyPath, local_location, fmt.Sprintf("%s@%s:%s", n.Username, n.Ip, remote_location))
	output, err := scpCmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("failed to upload file via %v, err:%v, output:%s for node:%v", fmt.Sprintf("%v\n", scpCmd), err, string(output), n.Id))
	} else {
		n.Logger.Debug(fmt.Sprintf("Upload using %v successful for node %v", fmt.Sprintf("%v\n", scpCmd), n.Id), 0)
	}
	return nil
}

// shut down the node

func (n *Node) Shut_Down() error {
	// shut down the node
	return n.ExecCmd("sudo shutdown -h now")
}

// start client

func (n *Node) Start_Client() error {
	// start the client program
	if n.Logger.DebugOn {
		return n.ExecCmd(fmt.Sprintf("cd %vbench/ && ./bench --node_config %vbench/ip.yaml --id %v --debug_on --debug_level %v", n.HomeDir, n.HomeDir, n.Id, n.Logger.Level))
	} else {
		return n.ExecCmd(fmt.Sprintf("cd %vbench/ && ./bench --node_config %vbench/ip.yaml --id %v", n.HomeDir, n.HomeDir, n.Id))
	}
}

func (n *Node) UpdateStats(perf []float64) {
	n.statMutex.Lock()
	n.stat.cpu_usage = perf[0]
	n.stat.mem_usage = perf[1]
	n.stat.network_in = perf[2]
	n.stat.network_out = perf[3]
	n.statMutex.Unlock()
}

func (n *Node) GetStats() NodeStat {
	n.statMutex.Lock()
	stats := NodeStat{
		cpu_usage:   n.stat.cpu_usage,
		mem_usage:   n.stat.mem_usage,
		network_in:  n.stat.network_in,
		network_out: n.stat.network_out,
	}
	n.statMutex.Unlock()
	return stats
}

func GetNodes(filename string) []*Node {
	print("reading node data from filename: \n", filename)
	type Nodes struct {
		Nodes []*Node `yaml:"nodes"`
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Error reading file: " + err.Error())
	}

	var nodes Nodes
	err = yaml.Unmarshal(data, &nodes)
	if err != nil {
		panic("Error unmarshalling YAML: " + err.Error())
	}

	client_nodes := nodes.Nodes[1:]
	// print the nodes
	print("Nodes: \n")
	for i := 0; i < len(client_nodes); i++ {
		fmt.Printf("Node: %v\n", client_nodes[i])
	}
	return client_nodes

}

func GetController(filename string) *Node {
	type Nodes struct {
		Nodes []*Node `yaml:"nodes"`
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Error reading file: " + err.Error())
	}

	var nodes Nodes
	err = yaml.Unmarshal(data, &nodes)
	if err != nil {
		panic("Error unmarshalling YAML: " + err.Error())
	}

	controller := nodes.Nodes[0]
	// print the controller

	fmt.Printf("Controller: %v\n", controller)
	return controller
}
