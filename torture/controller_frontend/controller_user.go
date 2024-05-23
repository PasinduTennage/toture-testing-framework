package controller_frontend

import (
	"time"
	"toture-test/torture/torture/src"
)

// main attack logic goes here

func StartAttack(nodes []torture.Attacker) {
	start_time := time.Now()
	for time.Now().Sub(start_time) < 60*time.Second {
		for _, node := range nodes {
			//node.DelayAllPacketsBy(10)
			//node.LossPercentagePackets(10)
			//node.DuplicatePercentagePackets(10)
			//node.ReorderPercentagePackets(10)
			//node.CorruptPercentagePackets(10)
			node.Halt()
			time.Sleep(5 * time.Second)
			node.CorruptDB()
			node.Kill()
			//node.BufferAllMessages()
			//node.AllowMessages(10)
		}
	}
	for _, node := range nodes {
		node.Exit()
	}
}
