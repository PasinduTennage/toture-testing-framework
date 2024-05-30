package attacker_impl

import (
	"fmt"
	"strconv"
	"strings"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
	torture "toture-test/torture/torture/src"
	"toture-test/torture/util"
)

// local NetEm attacker allows only one attack at a time

type LocalNetEmAttacker struct {
	name               int
	debugOn            bool
	debugLevel         int
	nextCommands       [][]string
	ports_under_attack []string // ports under attack
	process_id         string   // process under attack

	handle      string
	parent_band string
	prios       []int

	under_attack   bool
	current_attack int32

	c *torture.TortureClient
}

// NewLocalNetEmAttacker creates a new LocalNetEmAttacker

func NewLocalNetEmAttacker(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig, c *torture.TortureClient) *LocalNetEmAttacker {
	l := &LocalNetEmAttacker{
		name:           name,
		debugOn:        debugOn,
		debugLevel:     debugLevel,
		nextCommands:   [][]string{},
		under_attack:   false,
		current_attack: -1,
		c:              c,
	}

	v, ok := config.Options["ports"]
	if !ok || v == "NA" {
		panic("local netem attacker requires ports to be specified")
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

	l.Init(cgf)

	l.setNetEmVariables(cgf)

	return l
}

func (l *LocalNetEmAttacker) Init(cgf configuration.InstanceConfig) {
	// if this is the first attacker client, then initiate the root qdisc
	if strconv.Itoa(l.name) == cgf.Peers[0].Name {
		util.RunCommand("tc", []string{"filter", "del", "dev", "lo"})
		util.RunCommand("tc", []string{"qdisc", "del", "dev", "lo", "root"})
		util.RunCommand("tc", []string{"qdisc", "add", "dev", "lo", "root", "handle", "1:", "prio", "bands", strconv.Itoa(3 + len(cgf.Peers))})
	}
}

func (l *LocalNetEmAttacker) setNetEmVariables(cgf configuration.InstanceConfig) {
	index := -1
	for i, peer := range cgf.Peers {
		if peer.Name == strconv.Itoa(l.name) {
			index = i
			break
		}
	}
	if index == -1 {
		panic("could not find the peer in the configuration")
	}
	l.handle = strconv.Itoa((index + 1) * 10)
	l.parent_band = "1:" + strconv.Itoa(3+index)
	l.prios = []int{}
	for i := index*10 + 1; i < (index+1)*10; i++ { //assumes maximum number of ports to be 10 per consensus replica
		l.prios = append(l.prios, i)
	}
}

func (l *LocalNetEmAttacker) ExecuteLastCommands() error {
	var err error
	for i := 0; i < len(l.nextCommands); i++ {
		if len(l.nextCommands[i]) == 0 {
			continue
		}
		err = util.RunCommand(l.nextCommands[i][0], l.nextCommands[i][1:])
	}
	l.nextCommands = [][]string{}
	return err
}

func (l *LocalNetEmAttacker) applyHandleToEachPort() {
	i := 0
	for _, port := range l.ports_under_attack {
		util.RunCommand("tc", []string{"filter", "add", "dev", "lo", "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(l.prios[i]), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_band})
		l.nextCommands = append(l.nextCommands, []string{"tc", "filter", "del", "dev", "lo", "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(l.prios[i]), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_band})
		i++
	}
}

func (l *LocalNetEmAttacker) sendControllerMessage(m string) {
	l.c.SendControllerMessage(&proto.Message{
		StrParams: []string{m},
	})
}

func (l *LocalNetEmAttacker) DelayPackets(delay int, on bool) error {

	if !on {
		if !l.under_attack || l.current_attack != torture.NewOperationTypes().DelayPackets {
			l.sendControllerMessage("Failed to stop the delay attack because delay attack is currently not in progress")
			return nil
		} else {
			l.ExecuteLastCommands()
			l.under_attack = false
			l.current_attack = -1
			return nil
		}
	} else {
		if l.under_attack {
			l.sendControllerMessage("Failed to execute a delay attack because another attack is in progress")
			return nil
		} else {
			l.ExecuteLastCommands()
			err := util.RunCommand("tc", []string{"qdisc", "add", "dev", "lo", "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(delay) + "ms"})
			l.applyHandleToEachPort()
			l.nextCommands = append(l.nextCommands, []string{"tc", "qdisc", "del", "dev", "lo", "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(delay) + "ms"})
			l.under_attack = true
			l.current_attack = torture.NewOperationTypes().DelayPackets
			return err
		}
	}
}

func (l *LocalNetEmAttacker) LossPackets(int, bool) error {
	return nil
}

func (l *LocalNetEmAttacker) DuplicatePackets(int, bool) error {
	return nil
}

func (l *LocalNetEmAttacker) ReorderPackets(int, bool) error {
	return nil
}

func (l *LocalNetEmAttacker) CorruptPackets(int, bool) error {
	return nil
}

func (l *LocalNetEmAttacker) Pause(on bool) error {

	if !on {
		if !l.under_attack || l.current_attack != torture.NewOperationTypes().Pause {
			l.sendControllerMessage("Failed to stop the pause attack because pause attack is currently not in progress")
			return nil
		} else {
			l.ExecuteLastCommands()
			l.under_attack = false
			l.current_attack = -1
			return nil
		}
	} else {
		if l.under_attack {
			l.sendControllerMessage("Failed to execute a pause attack because another attack is in progress")
			return nil
		} else {
			l.ExecuteLastCommands()
			err := util.RunCommand("kill", []string{"-STOP", l.process_id})
			l.nextCommands = [][]string{{"kill", "-CONT", l.process_id}}
			l.under_attack = true
			l.current_attack = torture.NewOperationTypes().Pause
			return err
		}
	}
}

func (l *LocalNetEmAttacker) ResetAll() error {
	return l.ExecuteLastCommands()
}

func (l *LocalNetEmAttacker) Kill() error {
	l.ExecuteLastCommands()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	return err
}

func (l *LocalNetEmAttacker) QueueAllMessages(on bool) error {
	return nil
}

func (l *LocalNetEmAttacker) AllowMessages(int) error {
	return nil
}

func (l *LocalNetEmAttacker) CorruptDB() error {
	return nil
}

func (l *LocalNetEmAttacker) CleanUp() error {
	util.RunCommand("tc", []string{"filter", "del", "dev", "lo"})
	util.RunCommand("tc", []string{"qdisc", "del", "dev", "lo", "root"})
	return nil
}
