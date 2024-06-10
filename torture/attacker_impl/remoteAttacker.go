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

type RemoteNetEmAttacker struct {
	name               int
	debugOn            bool
	debugLevel         int
	nextNetEmCommands  [][]string // only for tc commands
	ports_under_attack []string   // ports under attack
	process_id         string     // process under attack

	device string

	handle      string
	parent_band string

	c *torture.TortureClient

	delayPackets     int
	lossPackets      int
	duplicatePackets int
	reorderPackets   int
	corruptPackets   int
}

// NewRemoteNetEmAttacker creates a new RemoteNetEmAttacker

func NewRemoteNetEmAttacker(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig, c *torture.TortureClient) *RemoteNetEmAttacker {
	l := &RemoteNetEmAttacker{
		name:              name,
		debugOn:           debugOn,
		debugLevel:        debugLevel,
		nextNetEmCommands: [][]string{},
		handle:            "20",
		parent_band:       "1:2",
		c:                 c,

		delayPackets:     0,
		lossPackets:      0,
		duplicatePackets: 0,
		reorderPackets:   0,
		corruptPackets:   0,
	}

	v, ok := config.Options["ports"]
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

	v, ok = config.Options["socket"]
	if !ok || v == "NA" {
		panic("remote netem attacker requires socket to be specified")
	} else {
		l.device = v
	}

	fmt.Printf("Process ID: %v, ports under attack %v, device %v\n", l.process_id, l.ports_under_attack, l.device)

	l.Init(cgf)

	return l
}

func (l *RemoteNetEmAttacker) Init(cgf configuration.InstanceConfig) {
	util.RunCommand("tc", []string{"filter", "del", "dev", l.device})
	util.RunCommand("tc", []string{"qdisc", "del", "dev", l.device, "root"})
	util.RunCommand("tc", []string{"qdisc", "add", "dev", l.device, "root", "handle", "1:", "prio", "bands", strconv.Itoa(5)})

}

func (l *RemoteNetEmAttacker) ExecuteLastNetEmCommands() error {
	var err error
	for i := 0; i < len(l.nextNetEmCommands); i++ {
		if len(l.nextNetEmCommands[i]) == 0 {
			continue
		}
		err = util.RunCommand(l.nextNetEmCommands[i][0], l.nextNetEmCommands[i][1:])
	}
	l.nextNetEmCommands = [][]string{}
	return err
}

func (l *RemoteNetEmAttacker) applyHandleToEachPort() {
	i := 1
	for _, port := range l.ports_under_attack {
		util.RunCommand("tc", []string{"filter", "add", "dev", l.device, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(i), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_band})
		l.nextNetEmCommands = append(l.nextNetEmCommands, []string{"tc", "filter", "del", "dev", l.device, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(i), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_band})
		i++
	}
}

func (l *RemoteNetEmAttacker) sendControllerMessage(m string) {
	l.c.SendControllerMessage(&proto.Message{
		StrParams: []string{m},
	})
}

func (l *RemoteNetEmAttacker) decrementDelay() {
	l.delayPackets--
}

func (l *RemoteNetEmAttacker) SetNewHandler() error {
	if l.reorderPackets > 0 && l.delayPackets == 0 {
		l.delayPackets = 1
		defer l.decrementDelay()
	}

	err := util.RunCommand("tc", []string{"qdisc", "add", "dev", l.device, "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(l.delayPackets) + "ms", "loss", strconv.Itoa(l.lossPackets) + "%", "duplicate", strconv.Itoa(l.duplicatePackets) + "%", "reorder", strconv.Itoa(l.reorderPackets) + "%", "50%", "corrupt", strconv.Itoa(l.corruptPackets) + "%"})
	l.applyHandleToEachPort()
	l.nextNetEmCommands = append(l.nextNetEmCommands, []string{"tc", "qdisc", "del", "dev", l.device, "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(l.delayPackets) + "ms", "loss", strconv.Itoa(l.lossPackets) + "%", "duplicate", strconv.Itoa(l.duplicatePackets) + "%", "reorder", strconv.Itoa(l.reorderPackets) + "%", "50%", "corrupt", strconv.Itoa(l.corruptPackets) + "%"})
	return err
}

func (l *RemoteNetEmAttacker) DelayPackets(delay int) error {
	l.ExecuteLastNetEmCommands()
	l.delayPackets = delay
	return l.SetNewHandler()

}

func (l *RemoteNetEmAttacker) LossPackets(loss int) error {
	l.ExecuteLastNetEmCommands()
	l.lossPackets = loss
	return l.SetNewHandler()

}

func (l *RemoteNetEmAttacker) DuplicatePackets(dup int) error {
	l.ExecuteLastNetEmCommands()
	l.duplicatePackets = dup
	return l.SetNewHandler()
}

func (l *RemoteNetEmAttacker) ReorderPackets(re int) error {
	l.ExecuteLastNetEmCommands()
	l.reorderPackets = re
	return l.SetNewHandler()
}

func (l *RemoteNetEmAttacker) CorruptPackets(corrupt int) error {
	l.ExecuteLastNetEmCommands()
	l.corruptPackets = corrupt
	return l.SetNewHandler()
}

func (l *RemoteNetEmAttacker) Pause(on bool) error {

	if on {
		err := util.RunCommand("kill", []string{"-STOP", l.process_id})
		return err
	} else {
		return util.RunCommand("kill", []string{"-CONT", l.process_id})
	}
}

func (l *RemoteNetEmAttacker) ResetAll() error {
	l.delayPackets = 0
	l.lossPackets = 0
	l.duplicatePackets = 0
	l.reorderPackets = 0
	l.corruptPackets = 0
	util.RunCommand("kill", []string{"-CONT", l.process_id})
	return l.ExecuteLastNetEmCommands()
}

func (l *RemoteNetEmAttacker) Kill() error {
	l.ExecuteLastNetEmCommands()
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
	util.RunCommand("tc", []string{"filter", "del", "dev", l.device})
	util.RunCommand("tc", []string{"qdisc", "del", "dev", l.device, "root"})
	return nil
}
