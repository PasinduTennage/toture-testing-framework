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

type LocalFilterAttacker struct {
	name       int
	debugOn    bool
	debugLevel int

	ports_under_attack []string // ports under attack
	process_id         string   // process under attack

	c *torture.TortureClient
}

// NewLocalFilterAttacker creates a new LocalFilterAttacker

func NewLocalFilterAttacker(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig, c *torture.TortureClient) *LocalFilterAttacker {
	l := &LocalFilterAttacker{
		name:       name,
		debugOn:    debugOn,
		debugLevel: debugLevel,
		c:          c,
	}

	v, ok := config.Options["ports"]
	if !ok || v == "NA" {
		panic("local netfilter attacker requires ports to be specified")
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

func (l *LocalFilterAttacker) Init(cgf configuration.InstanceConfig) {

}

func (l *LocalFilterAttacker) sendControllerMessage(m string) {
	l.c.SendControllerMessage(&proto.Message{
		StrParams: []string{m},
	})
}

func (l *LocalFilterAttacker) DelayPackets(delay int) error {

	return nil

}

func (l *LocalFilterAttacker) LossPackets(loss int) error {
	return nil

}

func (l *LocalFilterAttacker) DuplicatePackets(dup int) error {
	return nil
}

func (l *LocalFilterAttacker) ReorderPackets(re int) error {
	return nil
}

func (l *LocalFilterAttacker) CorruptPackets(corrupt int) error {
	return nil
}

func (l *LocalFilterAttacker) Pause(on bool) error {

	if on {
		err := util.RunCommand("kill", []string{"-STOP", l.process_id})
		return err
	} else {
		return util.RunCommand("kill", []string{"-CONT", l.process_id})
	}
}

func (l *LocalFilterAttacker) ResetAll() error {

	util.RunCommand("kill", []string{"-CONT", l.process_id})
	return nil
}

func (l *LocalFilterAttacker) Kill() error {
	l.CleanUp()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	return err
}

func (l *LocalFilterAttacker) QueueAllMessages(on bool) error {
	return nil
}

func (l *LocalFilterAttacker) AllowMessages(int) error {
	return nil
}

func (l *LocalFilterAttacker) CorruptDB() error {
	l.sendControllerMessage("CorruptDB is not supported by LocalFilterAttacker")
	return nil
}

func (l *LocalFilterAttacker) CleanUp() error {
	return nil
}
