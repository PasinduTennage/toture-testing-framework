package cmd

import (
	"os/exec"
	"strconv"
)

// LocalNetEmAttacker is a struct that implements the Attacker interface
type LocalNetEmAttacker struct {
	replicaName   string
	ports         map[int][]int // for each process id, the set of open ports
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

	process_port_map := make(map[int][]int)
	operations := make(map[int]string)

	for i := 0; i < len(ports); i++ {
		pid := GetProcessID(ports[i][0])
		if pid == -1 {
			panic("error getting process id")
		}
		process_port_map[pid] = ports[i]
		operations[pid] = ""
	}

	lNEA := LocalNetEmAttacker{
		replicaName:   replicaName,
		ports:         process_port_map,
		delay:         delay,
		lossRate:      lossRate,
		duplicateRate: duplicateRate,
		reorderRate:   reorderRate,
		corruptRate:   corruptRate,
		operations:    operations,
	}

	return &lNEA
}

func (lna *LocalNetEmAttacker) Start() error {
	// initialize the qdisc
	exec.Command("tc qdisc del dev lo root ; sudo tc filter del dev lo parent ffff:")
	exec.Command("tc qdisc add dev lo root handle 1: prio")
	return nil
}

func (lna *LocalNetEmAttacker) End() error {
	// delete all rules and filters
	exec.Command("tc qdisc del dev lo root ; sudo tc filter del dev lo parent ffff:")
	return nil
}

func (lna *LocalNetEmAttacker) ExecuteLastCommand(pId int) error {
	if len(lna.operations[pId]) > 0 {
		exec.Command(lna.operations[pId])
		lna.operations[pId] = ""
	}
	return nil
}

func (lna *LocalNetEmAttacker) Delay(pId int, delay int) error {
	lna.ExecuteLastCommand(pId)

	return nil

}

func (lna *LocalNetEmAttacker) Loss(pId int, lossRate int) error {
	lna.ExecuteLastCommand(pId)
	return nil
}

func (lna *LocalNetEmAttacker) Duplicate(pId int, duplicateRate int) error {
	lna.ExecuteLastCommand(pId)
	return nil
}

func (lna *LocalNetEmAttacker) Reorder(pId int, reorderRate int) error {
	lna.ExecuteLastCommand(pId)
	return nil
}

func (lna *LocalNetEmAttacker) Corrupt(pId int, corruptRate int) error {
	lna.ExecuteLastCommand(pId)
	return nil
}

func (lna *LocalNetEmAttacker) Halt(pId int) error {
	lna.ExecuteLastCommand(pId)
	exec.Command("kill -STOP " + strconv.Itoa(pId))
	lna.operations[pId] = "kill -CONT " + strconv.Itoa(pId)
	return nil
}

func (lna *LocalNetEmAttacker) Reset(pId int) error {
	lna.ExecuteLastCommand(pId)
	return nil
}

func (lna *LocalNetEmAttacker) kill(pId int) error {
	exec.Command("pkill -P " + strconv.Itoa(pId))
	return nil
}
