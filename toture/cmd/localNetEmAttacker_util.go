package cmd

import (
	"fmt"
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

	debugOn    bool
	debugLevel int
}

// NewLocalNetEmAttacker creates a new LocalNetEmAttacker

func NewLocalNetEmAttacker(
	replicaName string,
	ports [][]int,
	options map[string]any,
	debugOn bool, debugLevel int) *LocalNetEmAttacker {

	process_port_map := make(map[int][]int)
	operations := make(map[int]string)

	for i := 0; i < len(ports); i++ {
		pid := GetProcessID(ports[i][0])
		// fmt.Printf("pid: %v for port %v\n", pid, ports[i][0])
		if pid == -1 {
			panic("error getting process id")
		}
		process_port_map[pid] = ports[i]
		operations[pid] = ""
	}

	lNEA := LocalNetEmAttacker{
		replicaName:   replicaName,
		ports:         process_port_map,
		delay:         options["delay"].(int),
		lossRate:      options["lossRate"].(int),
		duplicateRate: options["duplicateRate"].(int),
		reorderRate:   options["reorderRate"].(int),
		corruptRate:   options["corruptRate"].(int),
		operations:    operations,
		debugOn:       debugOn,
		debugLevel:    debugLevel,
	}

	lNEA.debug("created a new LocalNetEmAttacker "+fmt.Sprintf("%v", lNEA), 1)

	return &lNEA
}

func (lna *LocalNetEmAttacker) Start() error {
	// initialize the qdisc
	exec.Command("tc qdisc del dev lo root ; sudo tc filter del dev lo parent ffff:")
	exec.Command("tc qdisc add dev lo root handle 1: prio")
	lna.debug("started LocalNetEmAttacker", 1)
	return nil
}

func (lna *LocalNetEmAttacker) End() error {
	// delete all rules and filters
	exec.Command("tc qdisc del dev lo root ; sudo tc filter del dev lo parent ffff:")
	lna.debug("ended LocalNetEmAttacker", 1)
	return nil
}

func (lna *LocalNetEmAttacker) ExecuteLastCommand(pId int) error {
	if len(lna.operations[pId]) > 0 {
		exec.Command(lna.operations[pId])
		lna.debug("executed last command for process "+strconv.Itoa(pId)+": "+lna.operations[pId], 1)
		lna.operations[pId] = ""
	}
	return nil
}

func (lna *LocalNetEmAttacker) Delay(pId int, delay int) error {
	lna.ExecuteLastCommand(pId)
	lna.debug("delaying process "+strconv.Itoa(pId)+" by "+strconv.Itoa(delay)+"ms", 1)

	return nil

}

func (lna *LocalNetEmAttacker) Loss(pId int, lossRate int) error {
	lna.debug("lossing process "+strconv.Itoa(pId)+" by "+strconv.Itoa(lossRate)+"%", 1)
	lna.ExecuteLastCommand(pId)
	return nil
}

func (lna *LocalNetEmAttacker) Duplicate(pId int, duplicateRate int) error {
	lna.ExecuteLastCommand(pId)
	lna.debug("duplicating process "+strconv.Itoa(pId)+" by "+strconv.Itoa(duplicateRate)+"%", 1)
	return nil
}

func (lna *LocalNetEmAttacker) Reorder(pId int, reorderRate int) error {
	lna.ExecuteLastCommand(pId)
	lna.debug("reordering process "+strconv.Itoa(pId)+" by "+strconv.Itoa(reorderRate)+"%", 1)
	return nil
}

func (lna *LocalNetEmAttacker) Corrupt(pId int, corruptRate int) error {
	lna.ExecuteLastCommand(pId)
	lna.debug("corrupting process "+strconv.Itoa(pId)+" by "+strconv.Itoa(corruptRate)+"%", 1)
	return nil
}

func (lna *LocalNetEmAttacker) Halt(pId int) error {
	lna.ExecuteLastCommand(pId)
	exec.Command("kill -STOP " + strconv.Itoa(pId))
	lna.debug("halting process "+strconv.Itoa(pId), 1)
	lna.operations[pId] = "kill -CONT " + strconv.Itoa(pId)
	return nil
}

func (lna *LocalNetEmAttacker) Reset(pId int) error {
	lna.ExecuteLastCommand(pId)
	lna.debug("resetting process "+strconv.Itoa(pId), 1)
	return nil
}

func (lna *LocalNetEmAttacker) Kill(pId int) error {
	exec.Command("pkill -P " + strconv.Itoa(pId))
	lna.debug("killing process "+strconv.Itoa(pId), 1)
	return nil
}

func (lna *LocalNetEmAttacker) GetPiDPortMap() map[int][]int {
	return lna.ports
}

func (lna *LocalNetEmAttacker) debug(s string, level int) {
	if lna.debugOn && lna.debugLevel >= level {
		println(s)
	}
}
