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
		links[0][1].SetDelay(200)
		links[0][2].SetDelay(300)
		links[0][3].SetDelay(400)
		links[0][4].SetDelay(500)

		time.Sleep(3 * time.Second)
		links[0][1].SetDelay(0)
		links[0][2].SetDelay(0)
		links[0][3].SetDelay(0)
		links[0][4].SetDelay(0)

		time.Sleep(1 * time.Second)
		//fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		//time.Sleep(1 * time.Second)
	}

	fmt.Print("Basic attack complete\n")
}
