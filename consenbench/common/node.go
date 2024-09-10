package common

import (
	"fmt"
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
	Id             int
	Ip             string
	Username       string
	HomeDir        string
	stat           NodeStat
	statMutex      *sync.Mutex
	privateKeyPath string
	Logger         *util.Logger
}

func NewNode(Id int, Ip string, Username string, HomeDir string, privateKeyPath string, logger *util.Logger) *Node {
	return &Node{
		Id:       Id,
		Ip:       Ip,
		Username: Username,
		HomeDir:  HomeDir,
		stat: NodeStat{
			cpu_usage:   0.0,
			mem_usage:   0.0,
			network_in:  0.0,
			network_out: 0.0,
		},
		statMutex:      &sync.Mutex{},
		privateKeyPath: privateKeyPath,
		Logger:         logger,
	}
}

// Execute a command on the node

func (n *Node) ExecCmd(cmd string) error {
	sshCmd := exec.Command("ssh", "-i", n.privateKeyPath, fmt.Sprintf("%s@%s", n.Username, n.Ip), cmd)
	output, err := sshCmd.CombinedOutput()
	if err != nil {
		n.Logger.Debug(fmt.Sprintf("failed to execute command via SSH, err:%v, output:%vs for node:%v", err, output, n.Id), 0)
		return fmt.Errorf("failed to execute command via SSH: %w, output: %s", err, output)
	} else {
		n.Logger.Debug(fmt.Sprintf("Executed command via SSH: %v, output: %v for node: %v", cmd, output, n.Id), 0)
	}
	return nil
}

// download the file from the remote location to the local location

func (n *Node) Get_Load(remote_location string, local_location string) error {

	scpCmd := exec.Command("scp", "-i", n.privateKeyPath, fmt.Sprintf("%s@%s:%s", n.Username, n.Ip, remote_location), local_location)
	output, err := scpCmd.CombinedOutput()
	if err != nil {
		n.Logger.Debug(fmt.Sprintf("failed to download file via SCP, errpr:%v, output:%v for node:%v", err, output, n.Id), 0)
		return fmt.Errorf("failed to download file via SCP: %w, output: %s", err, output)
	} else {
		n.Logger.Debug(fmt.Sprintf("Download using SCP successful for node %v", n.Id), 0)
	}
	return nil
}

// upload the file from the local location to the remote location

func (n *Node) Put_Load(local_location string, remote_location string) error {
	scpCmd := exec.Command("scp", "-i", n.privateKeyPath, local_location, fmt.Sprintf("%s@%s:%s", n.Username, n.Ip, remote_location))
	output, err := scpCmd.CombinedOutput()
	if err != nil {
		n.Logger.Debug(fmt.Sprintf("failed to upload file via SCP, err:%v, output:%s for node:%v", err, output, n.Id), 0)
		return fmt.Errorf("failed to upload file via SCP: %v, output: %s", err, output)
	} else {
		n.Logger.Debug(fmt.Sprintf("Upload using SCP successful for node %v", n.Id), 0)
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
	return n.ExecCmd(fmt.Sprintf("cd %vbench/ && ./bench --client --config %vbench/ip.yaml --name %v", n.HomeDir, n.HomeDir, n.Id))

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
