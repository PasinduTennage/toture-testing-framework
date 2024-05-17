package attacker_impl

import (
	"fmt"
	"strconv"
	"toture-test/torture/configuration"
	"toture-test/torture/util"
)

type LocalNetEmAttacker struct {
	name               int
	debugOn            bool
	debugLevel         int
	options            map[string]any
	nextCommand        string
	ports_under_attack []string // ports under attack
	process_id         string   // process under attack
}

// NewLocalNetEmAttacker creates a new LocalNetEmAttacker

func NewLocalNetEmAttacker(name int, debugOn bool, debugLevel int, options map[string]any, cgf configuration.InstanceConfig) *LocalNetEmAttacker {
	l := &LocalNetEmAttacker{
		name:        name,
		debugOn:     debugOn,
		debugLevel:  debugLevel,
		options:     options,
		nextCommand: "",
	}

	for i := 0; i < len(cgf.Peers); i++ {
		if cgf.Peers[i].Name == strconv.Itoa(name) {
			l.ports_under_attack = cgf.Peers[i].REPLICA_PORTS
			break
		}
	}
	if len(l.ports_under_attack) == 0 {
		panic("No ports to attack")
	}

	fmt.Printf("Ports under attack for client %v is: %v\n", name, l.ports_under_attack)

	l.process_id = strconv.Itoa(util.GetProcessID(l.ports_under_attack[0]))

	fmt.Printf("Process ID under attack: %v\n", l.process_id)

	return l
}

func (l *LocalNetEmAttacker) DelayAllPacketsBy(int) error {
	return nil
}

func (l *LocalNetEmAttacker) LossPercentagePackets(int) error {
	return nil
}

func (l *LocalNetEmAttacker) DuplicatePercentagePackets(int) error {
	return nil
}

func (l *LocalNetEmAttacker) ReorderPercentagePackets(int) error {
	return nil
}

func (l *LocalNetEmAttacker) CorruptPercentagePackets(int) error {
	return nil
}

func (l *LocalNetEmAttacker) Halt() error {
	return nil
}

func (l *LocalNetEmAttacker) Reset() error {
	return nil
}

func (l *LocalNetEmAttacker) Kill() error {
	return nil
}

func (l *LocalNetEmAttacker) BufferAllMessages() error {
	return nil
}

func (l *LocalNetEmAttacker) AllowMessages(int) error {
	return nil
}
