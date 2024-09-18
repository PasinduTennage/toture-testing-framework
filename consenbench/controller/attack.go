package controller

import (
	"strconv"
	"toture-test/consenbench/common"
	"toture-test/util"
)

func Max(i int, i2 int) int {
	if i > i2 {
		return i
	}
	return i2
}

func Sum(slice []float32) float32 {
	var total float32
	for _, val := range slice {
		total += val
	}
	return total
}

func GetAttackObjects(num_replicas int, replica_name string, nodes []*common.Node, controller *Controller, logger *util.Logger, ports []string) ([]*AttackNode, [][]*AttackLink, *LeaderOracle) {
	attackNodes := make([]*AttackNode, num_replicas)
	attackLinks := make([][]*AttackLink, num_replicas)
	leaderOracle := NewLeaderOracle(nodes[0:num_replicas], logger)

	id_ip_pairs := make([]string, num_replicas)
	for i := 0; i < num_replicas; i++ {
		id_ip_pairs[i] = strconv.Itoa(nodes[i].Id) + ":" + nodes[i].Ip
	}

	for i := 0; i < num_replicas; i++ {
		attackNodes[i] = NewAttackNode(nodes[i].Id, controller, replica_name, logger)
		attackNodes[i].Init(ports, id_ip_pairs)
		attackLinks[i] = make([]*AttackLink, num_replicas)
		for j := 0; j < num_replicas; j++ {
			if i != j {
				attackLinks[i][j] = NewAttackLink(nodes[i].Id, nodes[j].Id, controller, logger)
			}
		}
	}

	return attackNodes, attackLinks, leaderOracle
}

type Attack interface {
	Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int)
}
