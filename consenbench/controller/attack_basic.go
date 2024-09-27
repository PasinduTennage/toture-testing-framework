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
		links[0][1].SetBandwidth(1)
		links[0][2].SetBandwidth(2)

		//links[1][2].SetBandwidth(3)
		//links[1][3].SetBandwidth(4)
		//
		//links[2][3].SetBandwidth(5)
		//links[2][4].SetBandwidth(6)
		//
		//links[3][4].SetBandwidth(7)
		//links[3][0].SetBandwidth(8)
		//
		//links[4][0].SetBandwidth(9)
		//links[4][1].SetBandwidth(10)

		time.Sleep(3 * time.Second)

		links[0][1].SetBandwidth(10000000)
		links[0][2].SetBandwidth(10000000)

		//links[1][2].SetBandwidth(10000000)
		//links[1][3].SetBandwidth(10000000)
		//
		//links[2][3].SetBandwidth(10000000)
		//links[2][4].SetBandwidth(10000000)
		//
		//links[3][4].SetBandwidth(10000000)
		//links[3][0].SetBandwidth(10000000)
		//
		//links[4][0].SetBandwidth(10000000)
		//links[4][1].SetBandwidth(10000000)

		time.Sleep(1 * time.Second)
		//fmt.Printf("The leader order is %v\n", oracle.GetTopNLeaders())
		//time.Sleep(1 * time.Second)
	}

	fmt.Print("Basic attack complete\n")
}
