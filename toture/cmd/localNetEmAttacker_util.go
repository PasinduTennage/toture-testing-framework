package cmd

import "os/exec"

// LocalNetEmAttacker is a struct that implements the Attacker interface
type LocalNetEmAttacker struct {
	replicaName   string
	processIds    []int
	ports         [][]int // note that the order of ports and order of processIds are same
	delay         int
	lossRate      int
	duplicateRate int
	reorderRate   int
	corruptRate   int

	operations map[int]string // for each process, then next expected command
}

// NewLocalNetEmAttacker creates a new LocalNetEmAttacker

func NewLocalNetEmAttacker(
	replicaName string,
	ports [][]int,
	delay int,
	lossRate int,
	duplicateRate int,
	reorderRate int,
	corruptRate int) *LocalNetEmAttacker {

	p_ids := []int{}
	for i := 0; i < len(ports); i++ {
		pid := GetProcessID(ports[i][0])
		if pid == -1 {
			panic("error getting process id")
		}
		p_ids = append(p_ids, pid)
	}

	lNEA := LocalNetEmAttacker{
		replicaName:   replicaName,
		processIds:    p_ids,
		ports:         ports,
		delay:         delay,
		lossRate:      lossRate,
		duplicateRate: duplicateRate,
		reorderRate:   reorderRate,
		corruptRate:   corruptRate,
	}

	return &lNEA
}

func (lna *LocalNetEmAttacker) Start() error {
	// initialize the qdisc
	exec.Command("tc qdisc add dev lo root handle 1: prio")
	return nil
}

func (lna *LocalNetEmAttacker) End() error {
	// delete all rules and filters
	exec.Command("sudo tc qdisc del dev lo root ; sudo tc filter del dev lo parent ffff:")
	return nil
}

func (lna *LocalNetEmAttacker) Delay(pId int, delay int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetDelay(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Loss(pId int, lossRate int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetLoss(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Duplicate(pId int, duplicateRate int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetDuplicate(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Reorder(pId int, reorderRate int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetReorder(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Corrupt(pId int, corruptRate int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetCorrupt(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Halt(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetHalt(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) kill(pId int) error {
	return nil
}
