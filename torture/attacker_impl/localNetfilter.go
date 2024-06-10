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

type LocalNetFilterAttacker struct {
	name       int
	debugOn    bool
	debugLevel int

	ports_under_attack []string // ports under attack
	process_id         string   // process under attack

	c *torture.TortureClient
}

// NewLocalNetFilterAttacker creates a new LocalNetEmAttacker

func NewLocalNetFilterAttacker(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig, c *torture.TortureClient) *LocalNetFilterAttacker {
	l := &LocalNetFilterAttacker{
		name:       name,
		debugOn:    debugOn,
		debugLevel: debugLevel,
		c:          c,
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

func (l *LocalNetFilterAttacker) Init(cgf configuration.InstanceConfig) {

}

func (l *LocalNetFilterAttacker) sendControllerMessage(m string) {
	l.c.SendControllerMessage(&proto.Message{
		StrParams: []string{m},
	})
}

func (l *LocalNetFilterAttacker) DelayPackets(delay int) error {

	return nil

}

func (l *LocalNetFilterAttacker) LossPackets(loss int) error {
	return nil

}

func (l *LocalNetFilterAttacker) DuplicatePackets(dup int) error {
	return nil
}

func (l *LocalNetFilterAttacker) ReorderPackets(re int) error {
	return nil
}

func (l *LocalNetFilterAttacker) CorruptPackets(corrupt int) error {
	return nil
}

func (l *LocalNetFilterAttacker) Pause(on bool) error {

	if on {
		err := util.RunCommand("kill", []string{"-STOP", l.process_id})
		return err
	} else {
		return util.RunCommand("kill", []string{"-CONT", l.process_id})
	}
}

func (l *LocalNetFilterAttacker) ResetAll() error {

	util.RunCommand("kill", []string{"-CONT", l.process_id})
	return nil
}

func (l *LocalNetFilterAttacker) Kill() error {
	l.CleanUp()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	return err
}

func (l *LocalNetFilterAttacker) QueueAllMessages(on bool) error {
	return nil
}

func (l *LocalNetFilterAttacker) AllowMessages(int) error {
	return nil
}

func (l *LocalNetFilterAttacker) CorruptDB() error {
	l.sendControllerMessage("CorruptDB is not supported by LocalNetFilterAttacker")
	return nil
}

func (l *LocalNetFilterAttacker) CleanUp() error {
	return nil
}
