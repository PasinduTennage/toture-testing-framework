package util

import "sync"

type NodeStat struct {
	cpu_usage   float64
	mem_usage   float64
	network_in  float64
	network_out float64
}

type Node struct {
	Id        int
	Ip        string
	Username  string
	HomeDir   string
	stat      NodeStat
	statMutex sync.Mutex
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) ExecCmd(cmd string) error {
	// run the given command on the node
	return nil
}

func (n *Node) Get_Load(remote_location string, local_location string) error {
	// download the file from the remote location to the local location
	return nil
}

func (n *Node) Put_Load(local_location string, remote_location string) error {
	// upload the file from the local location to the remote location
	return nil
}

func (n *Node) Shut_Down() error {
	// shut down the node
	return nil
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
