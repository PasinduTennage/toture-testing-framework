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

func (n *Node) DelayPackets(delay int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().DelayPackets,
		IntParams: []int32{int32(delay)},
	})
	return nil
}

func (n *Node) LossPackets(loss_percentage int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().LossPackets,
		IntParams: []int32{int32(loss_percentage)},
	})
	return nil
}

func (n *Node) DuplicatePackets(duplicate_percentage int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().DuplicatePackets,
		IntParams: []int32{int32(duplicate_percentage)},
	})
	return nil
}

func (n *Node) ReorderPackets(reorder_percentage int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().ReorderPackets,
		IntParams: []int32{int32(reorder_percentage)},
	})
	return nil
}

func (n *Node) CorruptPackets(corrupt_percentage int) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().CorruptPackets,
		IntParams: []int32{int32(corrupt_percentage)},
	})
	return nil
}

func (n *Node) Pause(on bool) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().Pause,
		On:        on,
	})
	return nil
}

func (n *Node) ResetAll() error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().ResetAll,
	})
	return nil
}

func (n *Node) Kill() error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().Kill,
	})
	return nil
}

func (n *Node) QueueAllMessages(on bool) error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().QueueAllMessages,
		On:        on,
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

func (n *Node) CleanUp() error {
	n.backend.sendMessage(n.name, &proto.Message{
		Operation: NewOperationTypes().CleanUp,
	})
	return nil
}

func (c *TortureController) handleMessage(message *proto.Message, sender int) {
	// print message
	fmt.Printf("Controller received message %v from %d\n", message.StrParams[0], sender)
}
