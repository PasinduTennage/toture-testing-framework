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
	time.Sleep(time.Duration(duration) * time.Second)
	fmt.Print("Basic attack complete\n")
}
