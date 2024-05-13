package cmd

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
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
	cmd := exec.Command("tc", "qdisc", "del", "dev", "lo", "root")
	if cmd.Err != nil {
		panic(cmd.Err)
	}
	err := cmd.Run()
	if err != nil {
		println(cmd.Output())
	}

	cmd = exec.Command("tc", "filter", "del", "dev", "lo", "parent", "ffff:")
	if cmd.Err != nil {
		panic(cmd.Err)
	}
	err = cmd.Run()
	if err != nil {
		println(cmd.Output())
	}

	cmd = exec.Command("tc", "qdisc", "add", "dev", "lo", "root", "handle", "1:", "prio")
	if cmd.Err != nil {
		panic(cmd.Err)
	}
	err = cmd.Run()
	if err != nil {
		panic(err.Error())
	}
	lna.debug("started LocalNetEmAttacker", 1)
	return nil
}

func (lna *LocalNetEmAttacker) End() error {
	// delete all rules and filters
	cmd := exec.Command("tc", "qdisc", "del", "dev", "lo", "root")
	if cmd.Err != nil {
		panic(cmd.Err)
	}
	err := cmd.Run()
	if err != nil {
		println(cmd.Output())
		panic(err.Error())
	}

	cmd = exec.Command("tc", "filter", "del", "dev", "lo", "parent", "ffff:")
	if cmd.Err != nil {
		panic(cmd.Err)
	}
	err = cmd.Run()
	if err != nil {
		println(cmd.Output())
		//panic(err.Error())
	}

	lna.debug("ended LocalNetEmAttacker", 1)
	return nil
}

func (lna *LocalNetEmAttacker) ExecuteLastCommand(pId int) error {
	if len(lna.operations[pId]) > 0 {
		t := strings.Split(lna.operations[pId], " ")
		cmd := exec.Command(t[0], t[1:]...)
		if cmd.Err != nil {
			panic(cmd.Err)
		}
		err := cmd.Run()
		if err != nil {
			panic(err.Error())
		}
		lna.debug("executed last command for process "+strconv.Itoa(pId)+": "+lna.operations[pId], 1)
		lna.operations[pId] = ""
	}
	return nil
}

func (lna *LocalNetEmAttacker) Delay(pId int, delay int) error {
	lna.ExecuteLastCommand(pId)
	/*
	 tc qdisc del dev lo root
	 tc qdisc add dev lo root handle 1: prio
	 tc -s qdisc show dev lo
	 tc qdisc add dev lo parent 1:1 handle 10: netem delay 200ms
	 tc qdisc add dev lo parent 1:2 handle 20: netem delay 200ms
	 tc -s qdisc show dev lo
	 tc filter add dev lo protocol ip parent 1: prio 1 u32 match ip dport 10000 0xffff flowid 1:1 action flowid 1:10
	 tc filter add dev lo protocol ip parent 1: prio 1 u32 match ip dport 10001 0xffff flowid 1:1 action flowid 1:10
	 tc filter add dev lo protocol ip parent 1: prio 1 u32 match ip dport 20001 0xffff flowid 1:2 action flowid 1:20
	 tc filter add dev lo protocol ip parent 1: prio 1 u32 match ip dport 20002 0xffff flowid 1:2 action flowid 1:20
	 tc qdisc del dev lo parent 1:1 handle 10:
	 tc qdisc del dev lo parent 1:2 handle 20:
	*/

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
	cmd := exec.Command("kill", "-STOP", strconv.Itoa(pId))
	if cmd.Err != nil {
		panic(cmd.Err)
	}
	err := cmd.Run()
	if err != nil {
		panic(err.Error())
	}
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
	cmd := exec.Command("pkill", "-P", strconv.Itoa(pId))
	if cmd.Err != nil {
		panic(cmd.Err)
	}
	err := cmd.Run()
	if err != nil {
		panic(err.Error())
	}
	lna.debug("killing process "+strconv.Itoa(pId), 1)
	return nil
}

func (lna *LocalNetEmAttacker) GetPiDPortMap() map[int][]int {
	return lna.ports
}

func (lna *LocalNetEmAttacker) debug(s string, level int) {
	if lna.debugOn && lna.debugLevel <= level {
		println(s)
	}
}
