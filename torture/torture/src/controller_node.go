package torture

import (
	"fmt"
	"strconv"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
)

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

func CreateNodes(cfg configuration.InstanceConfig, backend *TortureController) []Attacker {
	nodes := make([]Attacker, 0)
	for i := 0; i < len(cfg.Peers); i++ {
		int_name, _ := strconv.Atoi(cfg.Peers[i].Name)
		node := NewNode(int64(int_name), backend)
		nodes = append(nodes, node)
	}
	return nodes
}

func (n *Node) DelayAllPacketsBy(delay int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().DelayAllPacketsBy,
		IntParams: []int32{int32(delay)},
	})
	return nil
}

func (n *Node) LossPercentagePackets(loss_percentage int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().LossPercentagePackets,
		IntParams: []int32{int32(loss_percentage)},
	})
	return nil
}

func (n *Node) DuplicatePercentagePackets(duplicate_percentage int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().DuplicatePercentagePackets,
		IntParams: []int32{int32(duplicate_percentage)},
	})
	return nil
}

func (n *Node) ReorderPercentagePackets(reorder_percentage int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().ReorderPercentagePackets,
		IntParams: []int32{int32(reorder_percentage)},
	})
	return nil
}

func (n *Node) CorruptPercentagePackets(corrupt_percentage int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().CorruptPercentagePackets,
		IntParams: []int32{int32(corrupt_percentage)},
	})
	return nil
}

func (n *Node) Halt() error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().Halt,
	})
	return nil
}

func (n *Node) Reset() error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().Reset,
	})
	return nil
}

func (n *Node) Kill() error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().Kill,
	})
	return nil
}

func (n *Node) BufferAllMessages() error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().BufferAllMessages,
	})
	return nil
}

func (n *Node) AllowMessages(num_messages int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().AllowMessages,
		IntParams: []int32{int32(num_messages)},
	})
	return nil
}

func (n *Node) CorruptDB() error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().CorruptDB,
	})
	return nil
}

func (n *Node) Exit() error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().Exit,
	})
	return nil
}

func (c *TortureController) handleMessage(message *proto.Message, sender int) {
	// print message
	c.debug(fmt.Sprintf("Controller received message %v from %d\n", message, sender), 0)
}
