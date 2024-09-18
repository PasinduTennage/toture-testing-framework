package controller

import (
	"fmt"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type LeaderOracle struct {
	nodes  []*common.Node
	logger *util.Logger
}

func NewLeaderOracle(nodes []*common.Node, logger *util.Logger) *LeaderOracle {
	return &LeaderOracle{
		nodes:  nodes,
		logger: logger,
	}
}

func (l *LeaderOracle) GetLeader() int {
	var leader *common.Node
	var highestCPU, highestNetIn, highestNetOut float32

	for _, node := range l.nodes {
		cpu_usage, _, network_in, network_out := node.GetStats()

		// Get the last 10 slots for each metric, or fewer if not enough data
		cpuUsage := cpu_usage[Max(0, len(cpu_usage)-3):]
		netIn := network_in[Max(0, len(network_in)-3):]
		netOut := network_out[Max(0, len(network_out)-3):]

		// Compute the sums of the last 10 slots
		cpuSum := Sum(cpuUsage)
		netInSum := Sum(netIn)
		netOutSum := Sum(netOut)

		// Check if this node has the highest values
		if cpuSum > highestCPU || netInSum > highestNetIn && netOutSum > highestNetOut {
			highestCPU = cpuSum
			highestNetIn = netInSum
			highestNetOut = netOutSum
			leader = node
		}
	}

	if leader == nil {
		panic("No leader found")
	}

	l.logger.Debug(fmt.Sprintf("Likely Leader is %v, with cpu: %v, net-in: %v, net-out: %v", leader.Id, highestCPU, highestNetIn, highestNetOut), 3)

	return leader.Id
}
