package cmd

import (
	"fmt"
	"time"
)

type Toture struct {
	attacker Attacker
	options  map[string]any
}

func NewTorture(attacker Attacker, options map[string]any) *Toture {
	to := Toture{
		attacker: attacker,
		options:  options,
	}
	return &to
}

// high level attack interface

func (t *Toture) Run() error {
	t.attacker.Start()
	t.startLocalSimpleAttack()
	t.attacker.End()
	return nil
}

// example attack that randomly slows down the replicas

func (t *Toture) startLocalSimpleAttack() {
	test_time, ok := t.options["testTime"]
	if !ok {
		panic("testTime option not found")
	}
	start := time.Now()
	for time.Now().Sub(start).Seconds() < float64(test_time.(int)) {
		// get two random pIds from lna.ports
		pids := t.getRandomProcessIDs(t.attacker.GetPiDPortMap(), t.options["numThreshold"].(int))
		for _, pid := range pids {
			t.attacker.Halt(pid)
		}
		time.Sleep(time.Duration(t.options["epochTime"].(int)) * time.Second)
		for _, pid := range pids {
			t.attacker.Reset(pid)
		}
	}
}

// get n random process ids from ports map
func (t *Toture) getRandomProcessIDs(ports map[int][]int, n int) []int {
	pids := []int{}
	for k, _ := range ports {
		pids = append(pids, k)
		if len(pids) == n {
			return pids
		}
	}
	panic("should not happen: pids: " + fmt.Sprintf("%v", pids) + ", port map: " + fmt.Sprintf("%v", ports) + ", n: " + fmt.Sprintf("%v", n))

}
