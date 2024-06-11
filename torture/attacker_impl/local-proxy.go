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

type Local_Proxy struct {
	name               int
	debugOn            bool
	debugLevel         int
	nextNetEmCommands  [][]string // only for tc commands
	ports_under_attack []string   // ports under attack
	process_id         string     // process under attack

	c *torture.TortureClient

	delayPackets     int
	lossPackets      int
	duplicatePackets int
	reorderPackets   int
	corruptPackets   int
}

// NewLocal_Proxy creates a new Local_Proxy

func NewLocal_Proxy(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig, c *torture.TortureClient) *Local_Proxy {
	l := &Local_Proxy{
		name:              name,
		debugOn:           debugOn,
		debugLevel:        debugLevel,
		nextNetEmCommands: [][]string{},
		c:                 c,

		delayPackets:     0,
		lossPackets:      0,
		duplicatePackets: 0,
		reorderPackets:   0,
		corruptPackets:   0,
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

	return l
}

func (l *Local_Proxy) Init(cgf configuration.InstanceConfig) {

	l.debug("initialized", 2)
}

func (l *Local_Proxy) sendControllerMessage(m string) {
	l.c.SendControllerMessage(&proto.Message{
		StrParams: []string{m},
	})
}

func (l *Local_Proxy) DelayPackets(delay int) error {
	l.delayPackets = delay
	l.debug("set new delay", 2)
	return nil

}

func (l *Local_Proxy) LossPackets(loss int) error {
	l.lossPackets = loss
	l.debug("set new loss", 2)
	return nil
}

func (l *Local_Proxy) DuplicatePackets(dup int) error {
	l.duplicatePackets = dup
	l.debug("set new duplication", 2)
	return nil
}

func (l *Local_Proxy) ReorderPackets(re int) error {
	l.reorderPackets = re
	l.debug("set new reorder", 2)
	return nil
}

func (l *Local_Proxy) CorruptPackets(corrupt int) error {
	l.corruptPackets = corrupt
	l.debug("set new corrupt", 2)
	return nil
}

func (l *Local_Proxy) Pause(on bool) error {

	if on {
		err := util.RunCommand("kill", []string{"-STOP", l.process_id})
		l.debug("paused", 2)
		return err
	} else {
		l.debug("resumed", 2)
		return util.RunCommand("kill", []string{"-CONT", l.process_id})
	}
}

func (l *Local_Proxy) ResetAll() error {
	l.delayPackets = 0
	l.lossPackets = 0
	l.duplicatePackets = 0
	l.reorderPackets = 0
	l.corruptPackets = 0
	util.RunCommand("kill", []string{"-CONT", l.process_id})
	l.debug("reset all", 2)
	return nil
}

func (l *Local_Proxy) Kill() error {
	l.CleanUp()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	l.debug("killed", 2)
	return err
}

func (l *Local_Proxy) QueueAllMessages(on bool) error {
	return nil
}

func (l *Local_Proxy) AllowMessages(n int) error {
	return nil
}

func (l *Local_Proxy) CorruptDB() error {
	l.sendControllerMessage("CorruptDB is not supported by Local_Proxy")
	return nil
}

func (l *Local_Proxy) CleanUp() error {
	return nil

}

func (l *Local_Proxy) debug(m string, level int) {
	if l.debugOn && level >= l.debugLevel {
		fmt.Println(m + "\n")
	}
}
