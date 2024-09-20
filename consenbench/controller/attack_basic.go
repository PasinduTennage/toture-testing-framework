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
		links[0][1].SetLoss(50)
		links[1][2].SetLoss(50)
		links[2][0].SetLoss(50)
		time.Sleep(1 * time.Second)
		links[0][1].SetLoss(0)
		links[1][2].SetLoss(0)
		links[2][0].SetLoss(0)
		time.Sleep(1 * time.Second)
		fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		time.Sleep(1 * time.Second)
	}

	fmt.Print("Basic attack complete\n")
}
