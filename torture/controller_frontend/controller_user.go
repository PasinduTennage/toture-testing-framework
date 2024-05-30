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
			node.DelayPackets(rand.Intn(20), true)
		}
		time.Sleep(2 * time.Second)
		for _, node := range nodes {
			node.DelayPackets(rand.Intn(20), true)
		}
		time.Sleep(2 * time.Second)
		for _, node := range nodes {
			node.Pause(true)
		}
		time.Sleep(2 * time.Second)
		for _, node := range nodes {
			node.Pause(false)
		}
		time.Sleep(2 * time.Second)
		for _, node := range nodes {
			node.DelayPackets(rand.Intn(20), false)
		}
	}
	for _, node := range nodes {
		node.Kill()
		node.CleanUp()
	}
	os.Exit(0)
}
