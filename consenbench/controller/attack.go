package controller

import (
	"toture-test/consenbench/common"
)

type AttackNode struct {
	Id           int
	Controller   *Controller
	Process_name string
}

func NewAttackNode(id int, controller *Controller, name string) *AttackNode {
	return &AttackNode{
		Id:           id,
		Controller:   controller,
		Process_name: name,
	}
}

func (n *AttackNode) Kill() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Kill),
				StringArgs:    []string{n.Process_name},
			},
		},
		Peer: n.Id,
	})
}

func (n *AttackNode) Slowdown() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Slowdown),
				StringArgs:    []string{n.Process_name},
			},
		},
		Peer: n.Id,
	})
}

func (n *AttackNode) Pause() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Pause),
				StringArgs:    []string{n.Process_name},
			},
		},
		Peer: n.Id,
	})
}

func (n *AttackNode) Continue() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Continue),
				StringArgs:    []string{n.Process_name},
			},
		},
		Peer: n.Id,
	})
}

func (n *AttackNode) SetSkew(ms float32) {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetSkew),
				FloatArgs:     []float32{ms},
			},
		},
		Peer: n.Id,
	})
}

func (n *AttackNode) SetDrift(ms float32) {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetDrift),
				FloatArgs:     []float32{ms},
			},
		},
		Peer: n.Id,
	})
}

type AttackLink struct {
	Id_sender   int
	Id_reciever int
	Controller  *Controller
}

func NewAttackLink(sender int, reciever int, controller *Controller) *AttackLink {
	return &AttackLink{
		Id_sender:   sender,
		Id_reciever: reciever,
		Controller:  controller,
	}
}

func (a *AttackLink) SetStatus(on bool) {
	if !on {
		a.SetLoss(100)
	} else {
		a.SetLoss(0)
	}
}

func (a *AttackLink) SetDelay(ms float32) {
	a.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetDelay),
				FloatArgs:     []float32{ms},
			},
		},
		Peer: a.Id_sender,
	})
}

func (a *AttackLink) SetLoss(l float32) {
	a.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetLoss),
				FloatArgs:     []float32{l},
			},
		},
		Peer: a.Id_sender,
	})
}

func (a *AttackLink) SetBandwidth(b float32) {
	a.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetBandwidth),
				FloatArgs:     []float32{b},
			},
		},
		Peer: a.Id_sender,
	})

}

type LeaderOracle struct {
	nodes []*common.Node
}

func NewLeaderOracle(nodes []*common.Node) *LeaderOracle {
	return &LeaderOracle{
		nodes: nodes,
	}
}

func (l *LeaderOracle) GetLeader() int {
	var leader *common.Node
	var highestCPU, highestNetIn, highestNetOut float32

	for _, node := range l.nodes {
		cpu_usage, _, network_in, network_out := node.GetStats()

		// Get the last 10 slots for each metric, or fewer if not enough data
		cpuUsage := cpu_usage[max(0, len(cpu_usage)-10):]
		netIn := network_in[max(0, len(network_in)-10):]
		netOut := network_out[max(0, len(network_out)-10):]

		// Compute the sums of the last 10 slots
		cpuSum := sum(cpuUsage)
		netInSum := sum(netIn)
		netOutSum := sum(netOut)

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
	return leader.Id
}

func max(i int, i2 int) int {
	if i > i2 {
		return i
	}
	return i2
}

func sum(slice []float32) float32 {
	var total float32
	for _, val := range slice {
		total += val
	}
	return total
}

func GetAttackObjects(num_replicas int, replica_name string, nodes []*common.Node, controller *Controller) ([]*AttackNode, [][]*AttackLink, *LeaderOracle) {
	attackNodes := make([]*AttackNode, num_replicas)
	attackLinks := make([][]*AttackLink, num_replicas)
	leaderOracle := NewLeaderOracle(nodes[0:num_replicas])

	for i := 0; i < num_replicas; i++ {
		attackNodes[i] = NewAttackNode(nodes[i].Id, controller, replica_name)
		attackLinks[i] = make([]*AttackLink, num_replicas)
		for j := 0; j < num_replicas; j++ {
			if i != j {
				attackLinks[i][j] = NewAttackLink(nodes[i].Id, nodes[j].Id, controller)
			}
		}
	}

	return attackNodes, attackLinks, leaderOracle
}

type Attack interface {
	Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int)
}
