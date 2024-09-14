package controller

import (
	"fmt"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type AttackNode struct {
	Id           int
	Controller   *Controller
	Process_name string
	logger       *util.Logger
}

func NewAttackNode(id int, controller *Controller, name string, logger *util.Logger) *AttackNode {
	return &AttackNode{
		Id:           id,
		Controller:   controller,
		Process_name: name,
		logger:       logger,
	}
}

func (n *AttackNode) Init(ports []string) {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().SetPorts),
				StringArgs:    append([]string{n.Process_name}, ports...),
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("Init node %v with process name %v and listening ports %v", n.Id, n.Process_name, ports), 3)
}

func (n *AttackNode) Kill() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Kill),
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("killed node %v", n.Id), 3)
}

func (n *AttackNode) Slowdown() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Slowdown),
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("Slowed node %v", n.Id), 3)
}

func (n *AttackNode) Pause() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Pause),
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("Paused node %v", n.Id), 3)
}

func (n *AttackNode) Continue() {
	n.Controller.Network.Send(&common.RPCPairPeer{
		RpcPair: &common.RPCPair{
			Code: common.GetRPCCodes().ControlMsg,
			Obj: &common.ControlMsg{
				OperationType: int32(common.GetOperationCodes().Continue),
			},
		},
		Peer: n.Id,
	})
	n.logger.Debug(fmt.Sprintf("Continued node %v", n.Id), 3)
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
	n.logger.Debug(fmt.Sprintf("Skewed node %v", n.Id), 3)
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
	n.logger.Debug(fmt.Sprintf("Drfted node %v", n.Id), 3)
}

type AttackLink struct {
	Id_sender   int
	Id_reciever int
	Controller  *Controller
	logger      *util.Logger
}

func NewAttackLink(sender int, reciever int, controller *Controller, logger *util.Logger) *AttackLink {
	return &AttackLink{
		Id_sender:   sender,
		Id_reciever: reciever,
		Controller:  controller,
		logger:      logger,
	}
}

func (a *AttackLink) SetStatus(on bool) {
	if !on {
		a.SetLoss(100)
	} else {
		a.SetLoss(0)
	}
	a.logger.Debug(fmt.Sprintf("set status link %v-%v", a.Id_sender, a.Id_reciever), 3)
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
	a.logger.Debug(fmt.Sprintf("set delay link %v-%v", a.Id_sender, a.Id_reciever), 3)
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
	a.logger.Debug(fmt.Sprintf("set loss link %v-%v", a.Id_sender, a.Id_reciever), 3)
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
	a.logger.Debug(fmt.Sprintf("set bandwidth link %v-%v", a.Id_sender, a.Id_reciever), 3)

}

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

	l.logger.Debug(fmt.Sprintf("Leader is %v", leader.Id), 3)

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

func GetAttackObjects(num_replicas int, replica_name string, nodes []*common.Node, controller *Controller, logger *util.Logger, ports []string) ([]*AttackNode, [][]*AttackLink, *LeaderOracle) {
	attackNodes := make([]*AttackNode, num_replicas)
	attackLinks := make([][]*AttackLink, num_replicas)
	leaderOracle := NewLeaderOracle(nodes[0:num_replicas], logger)

	for i := 0; i < num_replicas; i++ {
		attackNodes[i] = NewAttackNode(nodes[i].Id, controller, replica_name, logger)
		attackNodes[i].Init(ports)
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
