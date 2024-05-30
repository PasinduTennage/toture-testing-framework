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
	for time.Now().Sub(start_time) < 120*time.Second {
		for _, node := range nodes {
			node.LossPackets(rand.Intn(20), true)
		}
		time.Sleep(2 * time.Second)
		for _, node := range nodes {
			node.LossPackets(rand.Intn(20), false)
		}
		time.Sleep(2 * time.Second)
		for _, node := range nodes {
			node.DelayPackets(10, true)
		}
		time.Sleep(2 * time.Second)
		for _, node := range nodes {
			node.DelayPackets(10, false)
		}
		time.Sleep(2 * time.Second)
		for _, node := range nodes {
			node.ResetAll()
		}
	}
	for _, node := range nodes {
		node.CleanUp()
	}
	os.Exit(0)
}
