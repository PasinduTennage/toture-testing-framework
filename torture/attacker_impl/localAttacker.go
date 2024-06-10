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

type LocalNetEmAttacker struct {
	name               int
	debugOn            bool
	debugLevel         int
	nextNetEmCommands  [][]string // only for tc commands
	ports_under_attack []string   // ports under attack
	process_id         string     // process under attack

	handle      string
	parent_band string
	prios       []int

	c *torture.TortureClient

	delayPackets     int
	lossPackets      int
	duplicatePackets int
	reorderPackets   int
	corruptPackets   int

	packets <-chan netfilter.NFPacket
}

// NewLocalNetEmAttacker creates a new LocalNetEmAttacker

func NewLocalNetEmAttacker(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig, c *torture.TortureClient) *LocalNetEmAttacker {
	l := &LocalNetEmAttacker{
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

	l.setNetEmVariables(cgf)

	return l
}

func (l *LocalNetEmAttacker) Init(cgf configuration.InstanceConfig) {
	// if this is the first attacker client, then initiate the root qdisc
	if strconv.Itoa(l.name) == cgf.Peers[0].Name {
		util.RunCommand("tc", []string{"filter", "del", "dev", "lo"})
		util.RunCommand("tc", []string{"qdisc", "del", "dev", "lo", "root"})
		util.RunCommand("tc", []string{"qdisc", "add", "dev", "lo", "root", "handle", "1:", "prio", "bands", strconv.Itoa(3 + len(cgf.Peers))})
	}

	nfq, err := netfilter.NewNFQueue(uint16(l.name), 1000000, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		panic("could not open NFQUEUE: %v " + err.Error())
	}
	defer nfq.Close()

	packets := nfq.GetPackets()
	l.packets = packets

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
	for i := index*10 + 1; i < (index+1)*10; i++ { //assumes maximum number of ports to be 10 per consensus replica
		l.prios = append(l.prios, i)
	}
}

func (l *LocalNetEmAttacker) ExecuteLastNetEmCommands() error {
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

func (l *LocalNetEmAttacker) applyHandleToEachPort() {
	i := 0
	for _, port := range l.ports_under_attack {
		util.RunCommand("tc", []string{"filter", "add", "dev", "lo", "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(l.prios[i]), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_band})
		l.nextNetEmCommands = append(l.nextNetEmCommands, []string{"tc", "filter", "del", "dev", "lo", "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(l.prios[i]), "u32", "match", "ip", "dport", port, "0xffff", "flowid", l.parent_band})
		i++
	}
}

func (l *LocalNetEmAttacker) sendControllerMessage(m string) {
	l.c.SendControllerMessage(&proto.Message{
		StrParams: []string{m},
	})
}

func (l *LocalNetEmAttacker) decrementDelay() {
	l.delayPackets--
}

func (l *LocalNetEmAttacker) SetNewHandler() error {
	if l.reorderPackets > 0 && l.delayPackets == 0 {
		l.delayPackets = 1
		defer l.decrementDelay()
	}
	err := util.RunCommand("tc", []string{"qdisc", "add", "dev", "lo", "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(l.delayPackets) + "ms", "loss", strconv.Itoa(l.lossPackets) + "%", "duplicate", strconv.Itoa(l.duplicatePackets) + "%", "reorder", strconv.Itoa(l.reorderPackets) + "%", "50%", "corrupt", strconv.Itoa(l.corruptPackets) + "%"})
	l.applyHandleToEachPort()
	l.nextNetEmCommands = append(l.nextNetEmCommands, []string{"tc", "qdisc", "del", "dev", "lo", "parent", l.parent_band, "handle", l.handle + ":", "netem", "delay", strconv.Itoa(l.delayPackets) + "ms", "loss", strconv.Itoa(l.lossPackets) + "%", "duplicate", strconv.Itoa(l.duplicatePackets) + "%", "reorder", strconv.Itoa(l.reorderPackets) + "%", "50%", "corrupt", strconv.Itoa(l.corruptPackets) + "%"})
	return err
}

func (l *LocalNetEmAttacker) DelayPackets(delay int) error {
	l.ExecuteLastNetEmCommands()
	l.delayPackets = delay
	return l.SetNewHandler()

}

func (l *LocalNetEmAttacker) LossPackets(loss int) error {
	l.ExecuteLastNetEmCommands()
	l.lossPackets = loss
	return l.SetNewHandler()

}

func (l *LocalNetEmAttacker) DuplicatePackets(dup int) error {
	l.ExecuteLastNetEmCommands()
	l.duplicatePackets = dup
	return l.SetNewHandler()
}

func (l *LocalNetEmAttacker) ReorderPackets(re int) error {
	l.ExecuteLastNetEmCommands()
	l.reorderPackets = re
	return l.SetNewHandler()
}

func (l *LocalNetEmAttacker) CorruptPackets(corrupt int) error {
	l.ExecuteLastNetEmCommands()
	l.corruptPackets = corrupt
	return l.SetNewHandler()
}

func (l *LocalNetEmAttacker) Pause(on bool) error {

	if on {
		err := util.RunCommand("kill", []string{"-STOP", l.process_id})
		return err
	} else {
		return util.RunCommand("kill", []string{"-CONT", l.process_id})
	}
}

func (l *LocalNetEmAttacker) ResetAll() error {
	l.delayPackets = 0
	l.lossPackets = 0
	l.duplicatePackets = 0
	l.reorderPackets = 0
	l.corruptPackets = 0
	util.RunCommand("kill", []string{"-CONT", l.process_id})
	return l.ExecuteLastNetEmCommands()
}

func (l *LocalNetEmAttacker) Kill() error {
	l.ExecuteLastNetEmCommands()
	l.CleanUp()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	return err
}

func (l *LocalNetEmAttacker) QueueAllMessages(on bool) error {
	if on {
		// use iptables to redirect the traffic to ports to the queue with self.name
		for i := 0; i < len(l.ports_under_attack); i++ {
			port := l.ports_under_attack[i]
			util.RunCommand("iptables", []string{"-A", "INPUT", "-p", "tcp", "--dport", port, "-j", "NFQUEUE", "--queue-num", strconv.Itoa(l.name)})
		}
	} else {
		// use iptables to redirect the traffic to ports to the queue with self.name
		for i := 0; i < len(l.ports_under_attack); i++ {
			port := l.ports_under_attack[i]
			util.RunCommand("iptables", []string{"-D", "INPUT", "-p", "tcp", "--dport", port, "-j", "NFQUEUE", "--queue-num", strconv.Itoa(l.name)})
		}
	}
	return nil
}

func (l *LocalNetEmAttacker) AllowMessages(n int) error {
	go func() {
		for i := 0; i < n; i++ {
			packet := <-l.packets
			packet.SetVerdict(netfilter.NF_ACCEPT)
		}
	}()
	return nil
}

func (l *LocalNetEmAttacker) CorruptDB() error {
	l.sendControllerMessage("CorruptDB is not supported by LocalNetEmAttacker")
	return nil
}

func (l *LocalNetEmAttacker) CleanUp() error {
	util.RunCommand("tc", []string{"filter", "del", "dev", "lo"})
	util.RunCommand("tc", []string{"qdisc", "del", "dev", "lo", "root"})
	return l.QueueAllMessages(false)

}
