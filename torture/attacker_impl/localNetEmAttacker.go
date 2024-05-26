package attacker_impl

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"toture-test/torture/configuration"
	"toture-test/torture/util"
)

type LocalNetEmAttacker struct {
	name               int
	debugOn            bool
	debugLevel         int
	nextCommands       [][]string
	ports_under_attack []string // ports under attack
	process_id         string   // process under attack
	handle             string
	parent_band        string
	prios              []int
}

// NewLocalNetEmAttacker creates a new LocalNetEmAttacker

func NewLocalNetEmAttacker(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig) *LocalNetEmAttacker {
	l := &LocalNetEmAttacker{
		name:         name,
		debugOn:      debugOn,
		debugLevel:   debugLevel,
		nextCommands: [][]string{},
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
		l.ExecuteLastCommands()
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
	for i := index*10 + 1; i < (index+1)*10; i++ {
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

func (l *LocalNetEmAttacker) DelayAllPacketsBy(delay int) error {
	l.ExecuteLastCommands()
	err := util.RunCommand("tc", []string{"qdisc", "add", "dev", "lo", "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(delay) + "ms"})
	l.applyHandleToEachPort()
	l.nextCommands = append(l.nextCommands, []string{"tc", "qdisc", "del", "dev", "lo", "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(delay) + "ms"})
	return err
}

func (l *LocalNetEmAttacker) LossPercentagePackets(int) error {
	l.ExecuteLastCommands()
	return nil
}

func (l *LocalNetEmAttacker) DuplicatePercentagePackets(int) error {
	l.ExecuteLastCommands()
	return nil
}

func (l *LocalNetEmAttacker) ReorderPercentagePackets(int) error {
	l.ExecuteLastCommands()
	return nil
}

func (l *LocalNetEmAttacker) CorruptPercentagePackets(int) error {
	l.ExecuteLastCommands()
	return nil
}

func (l *LocalNetEmAttacker) Halt() error {
	l.ExecuteLastCommands()
	err := util.RunCommand("kill", []string{"-STOP", l.process_id})
	l.nextCommands = [][]string{{"kill", "-CONT", l.process_id}}
	return err
}

func (l *LocalNetEmAttacker) Reset() error {
	return l.ExecuteLastCommands()
}

func (l *LocalNetEmAttacker) Kill() error {
	l.ExecuteLastCommands()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	return err
}

func (l *LocalNetEmAttacker) BufferAllMessages() error {
	l.ExecuteLastCommands()
	return nil
}

func (l *LocalNetEmAttacker) AllowMessages(int) error {
	l.ExecuteLastCommands()
	return nil
}

func (l *LocalNetEmAttacker) CorruptDB() error {
	l.ExecuteLastCommands()
	return nil
}

func (l *LocalNetEmAttacker) Exit() error {
	l.ExecuteLastCommands()
	util.RunCommand("tc", []string{"filter", "del", "dev", "lo"})
	util.RunCommand("tc", []string{"qdisc", "del", "dev", "lo", "root"})
	os.Exit(0)
	return nil
}
