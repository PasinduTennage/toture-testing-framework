package controller_frontend

import (
	"time"
	"toture-test/torture/torture/src"
)

// main attack logic goes here

func StartAttack(nodes []*torture.Node) {
	start_time := time.Now()
	for time.Now().Sub(start_time) < 60*time.Second {
		for _, node := range nodes {
			//node.DelayAllPacketsBy(10)
			//node.LossPercentagePackets(10)
			//node.DuplicatePercentagePackets(10)
			//node.ReorderPercentagePackets(10)
			//node.CorruptPercentagePackets(10)
			//node.Halt()
			//node.Reset()
			node.Kill()
			//node.BufferAllMessages()
			//node.AllowMessages(10)
		}

	}
}
