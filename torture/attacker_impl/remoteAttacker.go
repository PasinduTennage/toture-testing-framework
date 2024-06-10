package attacker_impl

import (
	"fmt"
	"github.com/AkihiroSuda/go-netfilter-queue"
	"strconv"
	"strings"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
	torture "toture-test/torture/torture/src"
	"toture-test/torture/util"
)

type RemoteAttacker struct {
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

	packets chan netfilter.NFPacket
}

// NewRemoteAttacker creates a new RemoteAttacker

func NewRemoteAttacker(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig, c *torture.TortureClient) *RemoteAttacker {
	l := &RemoteAttacker{
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
		packets:          make(chan netfilter.NFPacket, 1000000),
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

func (l *RemoteAttacker) Init(cgf configuration.InstanceConfig) {
	util.RunCommand("tc", []string{"filter", "del", "dev", l.device})
	util.RunCommand("tc", []string{"qdisc", "del", "dev", l.device, "root"})
	util.RunCommand("iptables", []string{"-F"})
	util.RunCommand("tc", []string{"qdisc", "add", "dev", l.device, "root", "handle", "1:", "prio", "bands", strconv.Itoa(5)})

	go func() {
		nfq, err := netfilter.NewNFQueue(uint16(l.name), 1000000, netfilter.NF_DEFAULT_PACKET_SIZE)
		if err != nil {
			panic("could not open NFQUEUE: %v " + err.Error())
		}
		defer nfq.Close()

		packets := nfq.GetPackets()
		for packet := range packets {
			l.packets <- packet
		}
	}()
	l.debug("started the NF queue", 2)
}

func (l *RemoteAttacker) ExecuteLastNetEmCommands() error {
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

func (l *RemoteAttacker) applyHandleToEachPort() {
	i := 1
	for _, port := range l.ports_under_attack {
		util.RunCommand("tc", []string{"filter", "add", "dev", l.device, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(i), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_band})
		l.nextNetEmCommands = append(l.nextNetEmCommands, []string{"tc", "filter", "del", "dev", l.device, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(i), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_band})
		i++
	}
}

func (l *RemoteAttacker) sendControllerMessage(m string) {
	l.c.SendControllerMessage(&proto.Message{
		StrParams: []string{m},
	})
}

func (l *RemoteAttacker) decrementDelay() {
	l.delayPackets--
}

func (l *RemoteAttacker) SetNewHandler() error {
	if l.reorderPackets > 0 && l.delayPackets == 0 {
		l.delayPackets = 1
		defer l.decrementDelay()
	}

	err := util.RunCommand("tc", []string{"qdisc", "add", "dev", l.device, "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(l.delayPackets) + "ms", "loss", strconv.Itoa(l.lossPackets) + "%", "duplicate", strconv.Itoa(l.duplicatePackets) + "%", "reorder", strconv.Itoa(l.reorderPackets) + "%", "50%", "corrupt", strconv.Itoa(l.corruptPackets) + "%"})
	l.applyHandleToEachPort()
	l.nextNetEmCommands = append(l.nextNetEmCommands, []string{"tc", "qdisc", "del", "dev", l.device, "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(l.delayPackets) + "ms", "loss", strconv.Itoa(l.lossPackets) + "%", "duplicate", strconv.Itoa(l.duplicatePackets) + "%", "reorder", strconv.Itoa(l.reorderPackets) + "%", "50%", "corrupt", strconv.Itoa(l.corruptPackets) + "%"})
	l.debug("set new net em handler", 2)
	return err
}

func (l *RemoteAttacker) DelayPackets(delay int) error {
	l.ExecuteLastNetEmCommands()
	l.delayPackets = delay
	l.debug("set delay", 2)
	return l.SetNewHandler()

}

func (l *RemoteAttacker) LossPackets(loss int) error {
	l.ExecuteLastNetEmCommands()
	l.lossPackets = loss
	l.debug("set loss", 2)
	return l.SetNewHandler()

}

func (l *RemoteAttacker) DuplicatePackets(dup int) error {
	l.ExecuteLastNetEmCommands()
	l.duplicatePackets = dup
	l.debug("set duplicate", 2)
	return l.SetNewHandler()
}

func (l *RemoteAttacker) ReorderPackets(re int) error {
	l.ExecuteLastNetEmCommands()
	l.reorderPackets = re
	l.debug("set reorder", 2)
	return l.SetNewHandler()
}

func (l *RemoteAttacker) CorruptPackets(corrupt int) error {
	l.ExecuteLastNetEmCommands()
	l.corruptPackets = corrupt
	l.debug("set corrupt", 2)
	return l.SetNewHandler()
}

func (l *RemoteAttacker) Pause(on bool) error {

	if on {
		util.RunCommand("kill", []string{"-STOP", l.process_id})
		l.debug("paused", 2)
		return nil
	} else {
		util.RunCommand("kill", []string{"-CONT", l.process_id})
		l.debug("continue", 2)
		return nil
	}
}

func (l *RemoteAttacker) ResetAll() error {
	l.delayPackets = 0
	l.lossPackets = 0
	l.duplicatePackets = 0
	l.reorderPackets = 0
	l.corruptPackets = 0
	util.RunCommand("kill", []string{"-CONT", l.process_id})
	l.QueueAllMessages(false)
	l.debug("reset all", 2)
	return l.ExecuteLastNetEmCommands()
}

func (l *RemoteAttacker) Kill() error {
	l.ExecuteLastNetEmCommands()
	l.CleanUp()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	l.debug("killed", 2)
	return err
}

func (l *RemoteAttacker) QueueAllMessages(on bool) error {
	if on {
		// use iptables to redirect the traffic to ports to the queue with self.name
		for i := 0; i < len(l.ports_under_attack); i++ {
			port := l.ports_under_attack[i]
			util.RunCommand("iptables", []string{"-A", "INPUT", "-p", "tcp", "--dport", port, "-j", "NFQUEUE", "--queue-num", strconv.Itoa(l.name)})
		}
		l.debug("started queueing", 2)
	} else {
		for i := 0; i < len(l.ports_under_attack); i++ {
			port := l.ports_under_attack[i]
			util.RunCommand("iptables", []string{"-D", "INPUT", "-p", "tcp", "--dport", port, "-j", "NFQUEUE", "--queue-num", strconv.Itoa(l.name)})
		}
		l.debug("stopped queueing", 2)
	}
	return nil
}

func (l *RemoteAttacker) AllowMessages(n int) error {
	go func() {
		for i := 0; i < n; i++ {
			packet := <-l.packets
			packet.SetVerdict(netfilter.NF_ACCEPT)
		}
		l.debug("allowed messages", 2)
	}()
	return nil
}

func (l *RemoteAttacker) CorruptDB() error {
	l.sendControllerMessage("CorruptDB is not supported by RemoteAttacker")
	return nil
}

func (l *RemoteAttacker) CleanUp() error {
	util.RunCommand("tc", []string{"filter", "del", "dev", l.device})
	util.RunCommand("tc", []string{"qdisc", "del", "dev", l.device, "root"})
	util.RunCommand("iptables", []string{"-F"})
	return l.QueueAllMessages(false)
}

func (l *RemoteAttacker) debug(m string, level int) {
	if l.debugOn && level >= l.debugLevel {
		fmt.Println(m + "\n")
	}
}
