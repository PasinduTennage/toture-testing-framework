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

		links[1][2].SetDelay(400)
		links[1][3].SetDelay(500)

		links[2][3].SetDelay(600)
		links[2][4].SetDelay(700)

		links[3][4].SetDelay(800)
		links[3][0].SetDelay(900)

		links[4][0].SetDelay(1000)
		links[4][1].SetDelay(1100)

		time.Sleep(3 * time.Second)

		links[0][1].SetDelay(0)
		links[0][2].SetDelay(0)

		links[1][2].SetDelay(0)
		links[1][3].SetDelay(0)

		links[2][3].SetDelay(0)
		links[2][4].SetDelay(0)

		links[3][4].SetDelay(0)
		links[3][0].SetDelay(0)

		links[4][0].SetDelay(0)
		links[4][1].SetDelay(0)

		time.Sleep(1 * time.Second)
		//fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		//time.Sleep(1 * time.Second)
	}

	fmt.Print("Basic attack complete\n")
}
