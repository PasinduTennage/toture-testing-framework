package controller

import (
	"fmt"
	"time"
	"toture-test/util"
)

type BasicAttack struct {
	logger *util.Logger
}

func NewBasicAttack(logger *util.Logger) *BasicAttack {
	return &BasicAttack{
		logger: logger,
	}
}

func (a *BasicAttack) Attack(nodes []*AttackNode, links [][]*AttackLink, oracle *LeaderOracle, duration int) {
	fmt.Printf("Running basic attack for %v seconds\n", duration)
	start_time := time.Now()

	for time.Now().Sub(start_time).Seconds() < float64(duration-5) {
		for _, node := range nodes {
			node.Slowdown("true")
		}
		time.Sleep(1 * time.Second)
		for _, node := range nodes {
			node.Slowdown("false")
		}
		for i := 0; i < len(links); i++ {
			for j := 0; j < len(links[i]); j++ {
				if i == j {
					continue
				}
				links[i][j].SetLoss(10)
			}
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Print("Basic attack complete\n")
}
