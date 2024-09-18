package controller

import (
	"fmt"
	"time"
	"toture-test/consenbench/controller"
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

func (a *BasicAttack) Attack(nodes []*controller.AttackNode, links [][]*controller.AttackLink, oracle *controller.LeaderOracle, duration int) {
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
				links[i][j].SetLoss(50)
			}
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
		fmt.Printf("The leader is %v\n", oracle.GetLeader())
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
	}

	fmt.Print("Basic attack complete\n")
}
