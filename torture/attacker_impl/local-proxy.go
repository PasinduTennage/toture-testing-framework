package attacker_impl

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
	torture "toture-test/torture/torture/src"
	"toture-test/torture/util"
)

type Local_Proxy struct {
	name       int
	debugOn    bool
	debugLevel int

	ports_under_attack []string // ports under attack that the executor will be listening to
	dest_ports         []string
	dest_ip            string
	process_id         string // process under attack

	c *torture.TortureClient

	delayPackets         int
	lossPackets          int
	duplicatePackets     int
	reorderPackets       int
	corruptPackets       int
	paused               bool
	queued               bool
	allowMessageIfQueues chan bool

	mu *sync.RWMutex
}

// NewLocal_Proxy creates a new Local_Proxy

func NewLocal_Proxy(name int, debugOn bool, debugLevel int, cgf configuration.InstanceConfig, config configuration.ConsensusConfig, c *torture.TortureClient) *Local_Proxy {
	l := &Local_Proxy{
		name:       name,
		debugOn:    debugOn,
		debugLevel: debugLevel,
		c:          c,

		delayPackets:         0,
		lossPackets:          0,
		duplicatePackets:     0,
		reorderPackets:       0,
		corruptPackets:       0,
		paused:               false,
		queued:               false,
		allowMessageIfQueues: make(chan bool, 10000000),
		mu:                   &sync.RWMutex{},
	}

	v, ok := config.Options["ports"]
	if !ok || v == "NA" {
		panic("local proxy attacker requires ports to be specified")
	}
	l.ports_under_attack = strings.Split(v, " ")

	if len(l.ports_under_attack) == 0 {
		panic("no ports to attack")
	}

	v, ok = config.Options["dest_ports"]
	if !ok || v == "NA" {
		panic("local proxy attacker requires dest ports to be specified")
	}
	l.dest_ports = strings.Split(v, " ")

	if len(l.dest_ports) == 0 {
		panic("no dest ports to attack")
	}
	if len(l.dest_ports) != len(l.ports_under_attack) {
		panic("dest ports do not map the listening ports")
	}

	v, ok = config.Options["process_id"]
	if !ok || v == "NA" {
		l.process_id = strconv.Itoa(util.GetProcessID(l.dest_ports[0]))
		if l.process_id == "-1" {
			panic("could not find the process id")
		}
	} else {
		l.process_id = v
	}

	v, ok = config.Options["ip"]
	if !ok || v == "NA" {
		panic("could not find server ip")
	} else {
		l.dest_ip = v
	}

	fmt.Printf("Process ID: %v, ports under attack %v, destination ports: %v \n", l.process_id, l.ports_under_attack, l.dest_ports)

	l.Init(cgf)

	return l
}

func (l *Local_Proxy) Init(cgf configuration.InstanceConfig) {
	for i := 0; i < len(l.ports_under_attack); i++ {
		source_port := l.ports_under_attack[i]
		dest_port := l.dest_ports[i]
		l.runProxy(source_port, dest_port)
	}
	l.debug("initialized", 2)
}

func (l *Local_Proxy) runProxy(sPort string, dPort string) {
	go func(sPort string, dPort string) {
		// listen to the sPort in a new thread
		listener, err := net.Listen("tcp", "0.0.0.0:"+sPort)
		if err != nil {
			panic(err.Error())
		}
		l.debug("proxy listening to messages on "+"0.0.0.0:"+sPort, 0)
		for true {
			sConn, err := listener.Accept()
			if err != nil {
				panic(err.Error())
			}
			rConn, err := net.Dial("tcp", l.dest_ip+":"+dPort)
			if err != nil {
				panic(err.Error())
			}
			go l.runPipe(sConn, rConn)
		}
		// for each incoming new connection, open another thread that would process packets
	}(sPort, dPort)
}

func (l *Local_Proxy) runPipe(sCon net.Conn, dCon net.Conn) {
	for true {

	}
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
