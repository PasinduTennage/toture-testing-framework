package controller

import (
	"fmt"
	"sort"
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
		cpuUsage := cpu_usage[Max(0, len(cpu_usage)-5):]
		netIn := network_in[Max(0, len(network_in)-5):]
		netOut := network_out[Max(0, len(network_out)-5):]

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

// get the node ids of nodes which consume the most resources

func (l *LeaderOracle) GetTopNLeaders() []int {
	type NodeStats struct {
		node      *common.Node
		cpuSum    float32
		netInSum  float32
		netOutSum float32
	}

	var nodeStatsList []NodeStats

	// Gather stats for each node
	for _, node := range l.nodes {
		cpu_usage, _, network_in, network_out := node.GetStats()

		cpuUsage := cpu_usage[Max(0, len(cpu_usage)-5):]
		netIn := network_in[Max(0, len(network_in)-5):]
		netOut := network_out[Max(0, len(network_out)-5):]

		// Compute the sums of the last 3 slots
		cpuSum := Sum(cpuUsage)
		netInSum := Sum(netIn)
		netOutSum := Sum(netOut)

		nodeStatsList = append(nodeStatsList, NodeStats{
			node:      node,
			cpuSum:    cpuSum,
			netInSum:  netInSum,
			netOutSum: netOutSum,
		})
	}

	// Sort nodes by cpuSum first, then netInSum, and then netOutSum in descending order
	sort.Slice(nodeStatsList, func(i, j int) bool {
		if nodeStatsList[i].cpuSum != nodeStatsList[j].cpuSum {
			return nodeStatsList[i].cpuSum > nodeStatsList[j].cpuSum
		}
		if nodeStatsList[i].netInSum != nodeStatsList[j].netInSum {
			return nodeStatsList[i].netInSum > nodeStatsList[j].netInSum
		}
		return nodeStatsList[i].netOutSum > nodeStatsList[j].netOutSum
	})

	// Calculate how many leaders to select
	numLeaders := len(l.nodes)

	l.logger.Debug(fmt.Sprintf("Node resource usage --"), 3)

	// Extract the IDs of the top n leaders
	var leaderIds []int
	for i := 0; i < numLeaders && i < len(nodeStatsList); i++ {
		leaderIds = append(leaderIds, nodeStatsList[i].node.Id)
		l.logger.Debug(fmt.Sprintf("ID: %v, CPU: %v, NetIn: %v, NetOut: %v", nodeStatsList[i].node.Id, nodeStatsList[i].cpuSum, nodeStatsList[i].netInSum, nodeStatsList[i].netOutSum), 3)
	}

	return leaderIds
}
