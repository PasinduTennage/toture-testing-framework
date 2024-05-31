package attacker_impl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
	torture "toture-test/torture/torture/src"
	"toture-test/torture/util"
)

type RemoteNetEmAttacker struct {
	name       int
	debugOn    bool
	debugLevel int
	c          *torture.TortureClient

	ports_under_attack []string // ports under attack
	process_id         string   // process under attack

	handles      map[int32]string // a handle for each netem attack type
	parent_bands map[int32]string // parent band for each netem attack type
	prios        map[int32][]int  // for each netem attack type, a list of priorities

	attackValues map[int32]int        // for each netem attack type, the current value of the attack
	under_attack map[int32]bool       // for each attack type, whether it is currently in progress
	nextCommands map[int32][][]string // commands to run to stop each attack type

	networkAdapter string
}

// NewRemoteNetEmAttacker creates a new RemoteNetEmAttacker

func NewRemoteNetEmAttacker(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig, c *torture.TortureClient) *RemoteNetEmAttacker {
	l := &RemoteNetEmAttacker{
		name:       name,
		debugOn:    debugOn,
		debugLevel: debugLevel,
		c:          c,
	}

	v, ok := config.Options["adapter"]
	if !ok || v == "NA" {
		panic("remote netem adapter requires adapter to be specified")
	} else {
		l.networkAdapter = v
	}

	v, ok = config.Options["ports"]
	if !ok || v == "NA" {
		panic("remote netem attacker requires ports to be specified")
	}
	l.ports_under_attack = strings.Split(v, " ")

	if len(l.ports_under_attack) == 0 {
		panic("no ports to attack")
	}

	v, ok = config.Options["process_id"]
	if !ok || v == "NA" {
		l.process_id = strconv.Itoa(util.GetProcessID(l.ports_under_attack[0]))
	} else {
		l.process_id = v
	}
	fmt.Printf("Process ID: %v, ports under attack %v\n", l.process_id, l.ports_under_attack)

	l.Init()

	l.setNetEmVariables()

	return l
}

func (l *RemoteNetEmAttacker) Init() {
	util.RunCommand("tc", []string{"filter", "del", "dev", l.networkAdapter})
	util.RunCommand("tc", []string{"qdisc", "del", "dev", l.networkAdapter, "root"})
	util.RunCommand("tc", []string{"qdisc", "add", "dev", l.networkAdapter, "root", "handle", "1:", "prio", "bands", strconv.Itoa(4 + reflect.TypeOf(torture.OperationTypes{}).NumField())})

}

func (l *RemoteNetEmAttacker) setNetEmVariables() {
	l.handles = make(map[int32]string)
	l.parent_bands = make(map[int32]string)
	l.prios = make(map[int32][]int)
	l.attackValues = make(map[int32]int)
	l.under_attack = make(map[int32]bool)
	l.nextCommands = make(map[int32][][]string)

	l.handles[torture.NewOperationTypes().DelayPackets] = "10"
	l.handles[torture.NewOperationTypes().LossPackets] = "20"
	l.handles[torture.NewOperationTypes().DuplicatePackets] = "30"
	l.handles[torture.NewOperationTypes().ReorderPackets] = "40"
	l.handles[torture.NewOperationTypes().CorruptPackets] = "50"

	l.parent_bands[torture.NewOperationTypes().DelayPackets] = "1:3"
	l.parent_bands[torture.NewOperationTypes().LossPackets] = "1:4"
	l.parent_bands[torture.NewOperationTypes().DuplicatePackets] = "1:5"
	l.parent_bands[torture.NewOperationTypes().ReorderPackets] = "1:6"
	l.parent_bands[torture.NewOperationTypes().CorruptPackets] = "1:7"

	l.prios[torture.NewOperationTypes().DelayPackets] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	l.prios[torture.NewOperationTypes().LossPackets] = []int{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	l.prios[torture.NewOperationTypes().DuplicatePackets] = []int{21, 22, 23, 24, 25, 26, 27, 28, 29, 30}
	l.prios[torture.NewOperationTypes().ReorderPackets] = []int{31, 32, 33, 34, 35, 36, 37, 38, 39, 40}
	l.prios[torture.NewOperationTypes().CorruptPackets] = []int{41, 42, 43, 44, 45, 46, 47, 48, 49, 50}

	l.attackValues[torture.NewOperationTypes().DelayPackets] = 0
	l.attackValues[torture.NewOperationTypes().LossPackets] = 0
	l.attackValues[torture.NewOperationTypes().DuplicatePackets] = 0
	l.attackValues[torture.NewOperationTypes().ReorderPackets] = 0
	l.attackValues[torture.NewOperationTypes().CorruptPackets] = 0

	l.under_attack[torture.NewOperationTypes().DelayPackets] = false
	l.under_attack[torture.NewOperationTypes().LossPackets] = false
	l.under_attack[torture.NewOperationTypes().DuplicatePackets] = false
	l.under_attack[torture.NewOperationTypes().ReorderPackets] = false
	l.under_attack[torture.NewOperationTypes().CorruptPackets] = false
	l.under_attack[torture.NewOperationTypes().Pause] = false
	l.under_attack[torture.NewOperationTypes().ResetAll] = false
	l.under_attack[torture.NewOperationTypes().Kill] = false
	l.under_attack[torture.NewOperationTypes().QueueAllMessages] = false
	l.under_attack[torture.NewOperationTypes().AllowMessages] = false
	l.under_attack[torture.NewOperationTypes().CorruptDB] = false

	l.nextCommands[torture.NewOperationTypes().DelayPackets] = [][]string{}
	l.nextCommands[torture.NewOperationTypes().LossPackets] = [][]string{}
	l.nextCommands[torture.NewOperationTypes().DuplicatePackets] = [][]string{}
	l.nextCommands[torture.NewOperationTypes().ReorderPackets] = [][]string{}
	l.nextCommands[torture.NewOperationTypes().CorruptPackets] = [][]string{}
	l.nextCommands[torture.NewOperationTypes().Pause] = [][]string{}
	l.nextCommands[torture.NewOperationTypes().QueueAllMessages] = [][]string{}

}

func (l *RemoteNetEmAttacker) ExecuteLastCommands(aType int32) error {
	var err error
	for i := 0; i < len(l.nextCommands[aType]); i++ {
		if len(l.nextCommands[aType][i]) == 0 {
			continue
		}
		err = util.RunCommand(l.nextCommands[aType][i][0], l.nextCommands[aType][i][1:])
	}
	l.nextCommands[aType] = [][]string{}
	return err
}

func (l *RemoteNetEmAttacker) applyHandleToEachPort(aType int32) {
	i := 0
	for _, port := range l.ports_under_attack {
		util.RunCommand("tc", []string{"filter", "add", "dev", l.networkAdapter, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(l.prios[aType][i]), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_bands[aType]})
		l.nextCommands[aType] = append(l.nextCommands[aType], []string{"tc", "filter", "del", "dev", l.networkAdapter, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(l.prios[aType][i]), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_bands[aType]})
		i++
		if i == 10 {
			panic("only 10 maximum ports per consensus node supported")
		}
	}
}

func (l *RemoteNetEmAttacker) sendControllerMessage(m string) {
	l.c.SendControllerMessage(&proto.Message{
		StrParams: []string{m},
	})
}

func (l *RemoteNetEmAttacker) DelayPackets(delay int, on bool) error {

	if !on {
		if !l.under_attack[torture.NewOperationTypes().DelayPackets] {
			l.sendControllerMessage("Failed to stop the delay attack because delay attack is currently not in progress")
			return nil
		} else {
			l.ExecuteLastCommands(torture.NewOperationTypes().DelayPackets)
			l.under_attack[torture.NewOperationTypes().DelayPackets] = false
			l.attackValues[torture.NewOperationTypes().DelayPackets] = 0
			return nil
		}
	} else {
		if l.under_attack[torture.NewOperationTypes().DelayPackets] {
			l.ExecuteLastCommands(torture.NewOperationTypes().DelayPackets)
		}
		delay = delay + l.attackValues[torture.NewOperationTypes().DelayPackets]
		err := util.RunCommand("tc", []string{"qdisc", "add", "dev", l.networkAdapter, "parent", l.parent_bands[torture.NewOperationTypes().DelayPackets], "handle", l.handles[torture.NewOperationTypes().DelayPackets] + ":", "netem", "delay", strconv.Itoa(delay) + "ms"})
		l.applyHandleToEachPort(torture.NewOperationTypes().DelayPackets)
		l.nextCommands[torture.NewOperationTypes().DelayPackets] = append(l.nextCommands[torture.NewOperationTypes().DelayPackets], []string{"tc", "qdisc", "del", "dev", l.networkAdapter, "parent", l.parent_bands[torture.NewOperationTypes().DelayPackets], "handle", l.handles[torture.NewOperationTypes().DelayPackets] + ":", "netem", "delay", strconv.Itoa(delay) + "ms"})
		l.under_attack[torture.NewOperationTypes().DelayPackets] = true
		l.attackValues[torture.NewOperationTypes().DelayPackets] = delay
		return err

	}
}

func (l *RemoteNetEmAttacker) LossPackets(loss int, on bool) error {
	return nil
}

func (l *RemoteNetEmAttacker) DuplicatePackets(dup int, on bool) error {
	return nil
}

func (l *RemoteNetEmAttacker) ReorderPackets(re int, on bool) error {
	return nil

}

func (l *RemoteNetEmAttacker) CorruptPackets(corrupt int, on bool) error {
	return nil
}

func (l *RemoteNetEmAttacker) Pause(on bool) error {
	if !on {
		if !l.under_attack[torture.NewOperationTypes().Pause] {
			l.sendControllerMessage("Failed to stop the pause attack because pause attack is currently not in progress")
			return nil
		} else {
			l.ExecuteLastCommands(torture.NewOperationTypes().Pause)
			l.under_attack[torture.NewOperationTypes().Pause] = false
			return nil
		}
	} else {
		if l.under_attack[torture.NewOperationTypes().Pause] {
			l.sendControllerMessage("Pause attack already in progress")
			return nil
		} else {
			l.ExecuteLastCommands(torture.NewOperationTypes().Pause)
			err := util.RunCommand("kill", []string{"-STOP", l.process_id})
			l.nextCommands[torture.NewOperationTypes().Pause] = [][]string{{"kill", "-CONT", l.process_id}}
			l.under_attack[torture.NewOperationTypes().Pause] = true
			return err
		}
	}
}

func (l *RemoteNetEmAttacker) ResetAll() error {
	l.ExecuteLastCommands(torture.NewOperationTypes().DelayPackets)
	l.ExecuteLastCommands(torture.NewOperationTypes().LossPackets)
	l.ExecuteLastCommands(torture.NewOperationTypes().DuplicatePackets)
	l.ExecuteLastCommands(torture.NewOperationTypes().ReorderPackets)
	l.ExecuteLastCommands(torture.NewOperationTypes().CorruptPackets)
	l.ExecuteLastCommands(torture.NewOperationTypes().Pause)
	return l.ExecuteLastCommands(torture.NewOperationTypes().QueueAllMessages)
}

func (l *RemoteNetEmAttacker) Kill() error {
	l.ResetAll()
	l.CleanUp()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	return err
}

func (l *RemoteNetEmAttacker) QueueAllMessages(on bool) error {
	l.sendControllerMessage("QueueAllMessages is not supported by RemoteNetEmAttacker")
	return nil
}

func (l *RemoteNetEmAttacker) AllowMessages(int) error {
	l.sendControllerMessage("AllowMessages is not supported by RemoteNetEmAttacker")
	return nil
}

func (l *RemoteNetEmAttacker) CorruptDB() error {
	l.sendControllerMessage("CorruptDB is not supported by RemoteNetEmAttacker")
	return nil
}

func (l *RemoteNetEmAttacker) CleanUp() error {
	util.RunCommand("tc", []string{"filter", "del", "dev", l.networkAdapter})
	util.RunCommand("tc", []string{"qdisc", "del", "dev", l.networkAdapter, "root"})
	return nil
}
