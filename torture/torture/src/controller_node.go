package torture

import (
	"fmt"
	"strconv"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
)

type OperationTypes struct {
	DelayAllPacketsBy          int32
	LossPercentagePackets      int32
	DuplicatePercentagePackets int32
	ReorderPercentagePackets   int32
	CorruptPercentagePackets   int32
	Halt                       int32
	Reset                      int32
	Kill                       int32
	BufferAllMessages          int32
	AllowMessages              int32
}

func NewOperationTypes() OperationTypes {
	return OperationTypes{
		DelayAllPacketsBy:          1,
		LossPercentagePackets:      2,
		DuplicatePercentagePackets: 3,
		ReorderPercentagePackets:   4,
		CorruptPercentagePackets:   5,
		Halt:                       6,
		Reset:                      7,
		Kill:                       8,
		BufferAllMessages:          9,
		AllowMessages:              10,
	}
}

type Node struct {
	name    int64
	backend *TortureController
}

func NewNode(name int64, backend *TortureController) *Node {
	return &Node{
		name:    name,
		backend: backend,
	}
}

func CreateNodes(cfg configuration.InstanceConfig, backend *TortureController) []*Node {
	nodes := make([]*Node, 0)
	for i := 0; i < len(cfg.Peers); i++ {
		int_name, _ := strconv.Atoi(cfg.Peers[i].Name)
		node := NewNode(int64(int_name), backend)
		nodes = append(nodes, node)
	}
	return nodes
}

func (n *Node) DelayAllPacketsBy(delay int) {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().DelayAllPacketsBy,
		IntParams: []int32{int32(delay)},
	})
}

func (n *Node) LossPercentagePackets(loss_percentage int) {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().LossPercentagePackets,
		IntParams: []int32{int32(loss_percentage)},
	})
}

func (n *Node) DuplicatePercentagePackets(duplicate_percentage int) {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().DuplicatePercentagePackets,
		IntParams: []int32{int32(duplicate_percentage)},
	})
}

func (n *Node) ReorderPercentagePackets(reorder_percentage int) {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().ReorderPercentagePackets,
		IntParams: []int32{int32(reorder_percentage)},
	})
}

func (n *Node) CorruptPercentagePackets(corrupt_percentage int) {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().CorruptPercentagePackets,
		IntParams: []int32{int32(corrupt_percentage)},
	})
}

func (n *Node) Halt() {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().Halt,
	})
}

func (n *Node) Reset() {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().Reset,
	})
}

func (n *Node) Kill() {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().Kill,
	})
}

func (n *Node) BufferAllMessages() {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().BufferAllMessages,
	})
}

func (n *Node) AllowMessages(num_messages int) {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().AllowMessages,
		IntParams: []int32{int32(num_messages)},
	})
}

func (c *TortureController) handleMessage(message *proto.Message, sender int) {
	// print message
	c.debug(fmt.Sprintf("Controller received message %v from %d\n", message, sender), 0)
}
