package attacker_impl

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"toture-test/torture/configuration"
	"toture-test/torture/util"
)

type RemoteNetEmAttacker struct {
	name               int
	debugOn            bool
	debugLevel         int
	nextCommands       [][]string
	ports_under_attack []string // ports under attack
	process_id         string   // process under attack
	socket             string
}

// NewRemoteNetEmAttacker creates a new RemoteNetEmAttacker

func NewRemoteNetEmAttacker(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig) *RemoteNetEmAttacker {
	r := &RemoteNetEmAttacker{
		name:         name,
		debugOn:      debugOn,
		debugLevel:   debugLevel,
		nextCommands: [][]string{},
	}

	v, ok := config.Options["ports"]
	if !ok || v == "NA" {
		panic("remote netem attacker requires ports to be specified")
	}
	r.ports_under_attack = strings.Split(v, " ")

	if len(r.ports_under_attack) == 0 {
		panic("No ports to attack")
	}

	v, ok = config.Options["process_id"]
	if !ok || v == "NA" {
		r.process_id = strconv.Itoa(util.GetProcessID(r.ports_under_attack[0]))
	} else {
		r.process_id = v
	}

	v, ok = config.Options["socket"]
	if !ok || v == "NA" {
		panic("provide the network interface card name")
	} else {
		r.socket = v
	}

	fmt.Printf("Process ID %v, ports: %v, socket: %v\n", r.process_id, r.ports_under_attack, r.socket)

	return r
}

func (r *RemoteNetEmAttacker) ExecuteLastCommands() error {
	var err error
	for i := 0; i < len(r.nextCommands); i++ {
		if len(r.nextCommands) == 0 {
			continue
		}
		err = util.RunCommand(r.nextCommands[i][0], r.nextCommands[i][1:])
	}
	r.nextCommands = [][]string{}
	return err
}

func (r *RemoteNetEmAttacker) Init() error {
	r.ExecuteLastCommands()
	err := util.RunCommand("tc", []string{"filter", "del", "dev", r.socket})
	err = util.RunCommand("tc", []string{"qdisc", "del", "dev", r.socket, "root"})
	err = util.RunCommand("tc", []string{"qdisc", "add", "dev", r.socket, "root", "handle", "1:", "prio"})
	return err
}

func (r *RemoteNetEmAttacker) applyHandleToEachPort() {
	i := 1
	for _, port := range r.ports_under_attack {
		util.RunCommand("tc", []string{"filter", "add", "dev", r.socket, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(i), "u32", "match", "ip", "dport", port, "0xffff", "flowid", "1:1"})
		r.nextCommands = append(r.nextCommands, []string{"tc", "filter", "del", "dev", r.socket, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(i), "u32", "match", "ip", "dport", port, "0xffff", "flowid", "1:1"})
		i++
	}
}

func (r *RemoteNetEmAttacker) DelayAllPacketsBy(delay int) error {
	r.ExecuteLastCommands()
	err := util.RunCommand("tc", []string{"qdisc", "add", "dev", r.socket, "parent", "1:1", "handle", "10:", "netem", "delay", strconv.Itoa(delay) + "ms"})
	r.applyHandleToEachPort()
	r.nextCommands = append(r.nextCommands, []string{"tc", "qdisc", "del", "dev", r.socket, "parent", "1:1", "handle", "10:", "netem", "delay", strconv.Itoa(delay) + "ms"})
	return err
}

func (r *RemoteNetEmAttacker) LossPercentagePackets(int) error {
	r.ExecuteLastCommands()
	return nil
}

func (r *RemoteNetEmAttacker) DuplicatePercentagePackets(int) error {
	r.ExecuteLastCommands()
	return nil
}

func (r *RemoteNetEmAttacker) ReorderPercentagePackets(int) error {
	r.ExecuteLastCommands()
	return nil
}

func (r *RemoteNetEmAttacker) CorruptPercentagePackets(int) error {
	r.ExecuteLastCommands()
	return nil
}

func (r *RemoteNetEmAttacker) Halt() error {
	r.ExecuteLastCommands()
	err := util.RunCommand("kill", []string{"-STOP", r.process_id})
	r.nextCommands = [][]string{{"kill", "-CONT", r.process_id}}
	return err
}

func (r *RemoteNetEmAttacker) Reset() error {
	return r.ExecuteLastCommands()
}

func (r *RemoteNetEmAttacker) Kill() error {
	r.ExecuteLastCommands()
	err := util.RunCommand("kill", []string{"-9", r.process_id})
	return err
}

func (r *RemoteNetEmAttacker) BufferAllMessages() error {
	r.ExecuteLastCommands()
	return nil
}

func (r *RemoteNetEmAttacker) AllowMessages(int) error {
	r.ExecuteLastCommands()
	return nil
}

func (r *RemoteNetEmAttacker) CorruptDB() error {
	r.ExecuteLastCommands()
	return nil
}

func (r *RemoteNetEmAttacker) Exit() error {
	r.ExecuteLastCommands()
	os.Exit(0)
	return nil
}
