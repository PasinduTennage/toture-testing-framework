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
			node.DelayPackets(10, true)
		}
	}
}
