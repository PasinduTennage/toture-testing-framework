package controller

import (
	"fmt"
	"time"
)

type BasicAttack struct {
}

func NewBasicAttack() *BasicAttack {
	return &BasicAttack{}
}

func (a *BasicAttack) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running basic attack for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-5) {
		for _, node := range nodes {
			node.Slowdown()
		}
		time.Sleep(1 * time.Second)
		for i := 0; i < len(links); i++ {
			for j := 0; j < len(links[i]); j++ {
				if i == j {
					continue
				}
				links[i][j].SetLoss(0)
			}
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Print("Basic attack complete\n")
}
