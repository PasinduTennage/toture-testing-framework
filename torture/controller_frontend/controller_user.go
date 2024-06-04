package controller_frontend

import (
	"math/rand"
	"os"
	"time"
	"toture-test/torture/torture/src"
)

// main attack logic goes here

func StartAttack(nodes []torture.Attacker) {
	start_time := time.Now()
	for time.Now().Sub(start_time) < 20*time.Second {
		//delay
		for _, node := range nodes {
			node.DelayPackets(rand.Intn(20))
		}
		time.Sleep(2 * time.Second)
		//for _, node := range nodes {
		//	node.DelayPackets(0)
		//}
		//time.Sleep(2 * time.Second)
		//loss
		for _, node := range nodes {
			node.LossPackets(rand.Intn(20))
		}
		time.Sleep(2 * time.Second)
		//for _, node := range nodes {
		//	node.LossPackets(0)
		//}
		//time.Sleep(2 * time.Second)

		// duplicate
		for _, node := range nodes {
			node.DuplicatePackets(90)
		}
		time.Sleep(2 * time.Second)
		//for _, node := range nodes {
		//	node.DuplicatePackets(0)
		//}
		//time.Sleep(2 * time.Second)

		//reorder
		for _, node := range nodes {
			node.ReorderPackets(rand.Intn(20))
		}
		time.Sleep(2 * time.Second)
		//for _, node := range nodes {
		//	node.ReorderPackets(0)
		//}
		//time.Sleep(2 * time.Second)

		// corrupt
		for _, node := range nodes {
			node.CorruptPackets(20)
		}
		time.Sleep(2 * time.Second)
		//for _, node := range nodes {
		//	node.CorruptPackets(0)
		//}
		//time.Sleep(2 * time.Second)

		//pause
		for _, node := range nodes {
			node.Pause(true)
		}
		time.Sleep(2 * time.Second)
		//for _, node := range nodes {
		//	node.Pause(false)
		//}
		//time.Sleep(2 * time.Second)

		//resetAll
		for _, node := range nodes {
			node.ResetAll()
		}
		time.Sleep(2 * time.Second)

	}
	for _, node := range nodes {
		node.CleanUp()
	}
	os.Exit(0)
}
