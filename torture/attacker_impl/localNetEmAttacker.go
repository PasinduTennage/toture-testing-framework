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
		panic("No ports to attack")
	}

	v, ok = config.Options["process_id"]
	if !ok || v == "NA" {
		l.process_id = strconv.Itoa(util.GetProcessID(l.ports_under_attack[0]))
	} else {
		l.process_id = v
	}
	fmt.Printf("Process ID: %v, ports under attack %v\n", l.process_id, l.ports_under_attack)
	return l
}

func (l *LocalNetEmAttacker) ExecuteLastCommands() error {
	var err error
	for i := 0; i < len(l.nextCommands); i++ {
		if len(l.nextCommands[i]) == 0 {
			return nil
		}
		err = util.RunCommand(l.nextCommands[i][0], l.nextCommands[i][1:])
	}
	l.nextCommands = [][]string{}
	return err
}

func (l *LocalNetEmAttacker) DelayAllPacketsBy(int) error {
	l.ExecuteLastCommands()
	return nil
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
	os.Exit(0)
	return nil
}
