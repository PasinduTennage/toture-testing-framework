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
	options            map[string]any
	nextCommand        []string
	ports_under_attack []string // ports under attack
	process_id         string   // process under attack
}

// NewLocalNetEmAttacker creates a new LocalNetEmAttacker

func NewLocalNetEmAttacker(name int, debugOn bool, debugLevel int, options map[string]any, cgf configuration.InstanceConfig, config configuration.ConsensusConfig) *LocalNetEmAttacker {
	l := &LocalNetEmAttacker{
		name:        name,
		debugOn:     debugOn,
		debugLevel:  debugLevel,
		options:     options,
		nextCommand: []string{},
	}

	v, ok := config.Options["ports"]
	if !ok || v == "NA" {
		panic("local netem attacker requires ports to be specified")
	}
	l.ports_under_attack = strings.Split(v, " ")

	if len(l.ports_under_attack) == 0 {
		panic("No ports to attack")
	}

	fmt.Printf("Ports under attack for NetEm client %v is: %v\n", name, l.ports_under_attack)

	v, ok = config.Options["process_id"]
	if !ok || v == "NA" {
		l.process_id = strconv.Itoa(util.GetProcessID(l.ports_under_attack[0]))
	} else {
		l.process_id = v
	}

	fmt.Printf("Process ID under attack: %v\n", l.process_id)

	return l
}

func (l *LocalNetEmAttacker) DelayAllPacketsBy(int) error {
	l.ExecuteLastCommand()
	return nil
}

func (l *LocalNetEmAttacker) LossPercentagePackets(int) error {
	l.ExecuteLastCommand()
	return nil
}

func (l *LocalNetEmAttacker) DuplicatePercentagePackets(int) error {
	l.ExecuteLastCommand()
	return nil
}

func (l *LocalNetEmAttacker) ReorderPercentagePackets(int) error {
	l.ExecuteLastCommand()
	return nil
}

func (l *LocalNetEmAttacker) CorruptPercentagePackets(int) error {
	l.ExecuteLastCommand()
	return nil
}

func (l *LocalNetEmAttacker) Halt() error {
	l.ExecuteLastCommand()
	err := util.RunCommand("kill", []string{"-STOP", l.process_id})
	l.nextCommand = []string{"kill", "-CONT", l.process_id}
	return err
}

func (l *LocalNetEmAttacker) Reset() error {
	return l.ExecuteLastCommand()
}

func (l *LocalNetEmAttacker) Kill() error {
	l.ExecuteLastCommand()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	return err
}

func (l *LocalNetEmAttacker) BufferAllMessages() error {
	l.ExecuteLastCommand()
	return nil
}

func (l *LocalNetEmAttacker) AllowMessages(int) error {
	l.ExecuteLastCommand()
	return nil
}

func (l *LocalNetEmAttacker) CorruptDB() error {
	l.ExecuteLastCommand()
	return nil
}

func (l *LocalNetEmAttacker) Exit() error {
	l.ExecuteLastCommand()
	os.Exit(0)
	return nil
}

func (l *LocalNetEmAttacker) ExecuteLastCommand() error {
	if len(l.nextCommand) == 0 {
		return nil
	}
	err := util.RunCommand(l.nextCommand[0], l.nextCommand[1:])
	l.nextCommand = []string{}
	return err
}
