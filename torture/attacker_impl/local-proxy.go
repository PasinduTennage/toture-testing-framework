package attacker_impl

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
	torture "toture-test/torture/torture/src"
	"toture-test/torture/util"
)

type Local_Proxy struct {
	name       int
	debugOn    bool
	debugLevel int

	c *torture.TortureClient

	delayPackets     int
	lossPackets      int
	duplicatePackets int
	reorderPackets   int
	corruptPackets   int

	paused bool
	queued bool

	mu *sync.RWMutex // for mutual exclusion of above 7 variables

	allowMessageIfQueued chan bool // used to inform each connection handler whether to release the message or not

	listening_ports []string // ports under attack that the executor will be listening to
	dest_ports      []string // ports to which the consensus protocol is listening to
	dest_ip         string   // ip of the consensus replica
	process_id      string   // process id of the consensus replica

	numConnections int
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
		mu:                   &sync.RWMutex{},
		allowMessageIfQueued: make(chan bool, 1000000),
		numConnections:       0,
	}

	v, ok := config.Options["ports"]
	if !ok || v == "NA" {
		panic("local proxy attacker requires ports to be specified")
	}
	l.listening_ports = strings.Split(v, " ")

	if len(l.listening_ports) == 0 {
		panic("no listening ports specified")
	}

	v, ok = config.Options["dest_ports"]
	if !ok || v == "NA" {
		panic("local proxy attacker requires dest ports to be specified")
	}
	l.dest_ports = strings.Split(v, " ")

	if len(l.dest_ports) == 0 {
		panic("no dest ports to attack")
	}
	if len(l.dest_ports) != len(l.listening_ports) {
		panic("dest ports do not map the listening ports")
	}

	v, ok = config.Options["ip"]
	if !ok || v == "NA" {
		panic("could not find server ip")
	} else {
		l.dest_ip = v
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

	fmt.Printf("Local proxy attacher: Process ID: %v, listening ports %v, destination ports: %v, ip:%v \n", l.process_id, l.listening_ports, l.dest_ports, l.dest_ip)

	l.Init(cgf)

	return l
}

func (l *Local_Proxy) Init(cgf configuration.InstanceConfig) {
	for i := 0; i < len(l.listening_ports); i++ {
		source_port := l.listening_ports[i]
		dest_port := l.dest_ports[i]
		l.runProxy(source_port, dest_port)
	}
	l.debug("initialized all listening ports", 2)
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
			l.debug("new connection from remote replica to "+sPort, 2)
			rConn, err := net.Dial("tcp", l.dest_ip+":"+dPort)
			if err != nil {
				panic(err.Error())
			}
			l.debug("setup new connection to local replica "+dPort, 2)
			go l.runPipe(sConn, rConn)
		}
	}(sPort, dPort)
}

func (l *Local_Proxy) runPipe(sCon net.Conn, dCon net.Conn) {
	l.mu.Lock()
	l.numConnections++
	l.mu.Unlock()
	incoming := make(chan []byte, 1000)
	go func() {
		for true {
			buffer := make([]byte, 100)
			n, err := sCon.Read(buffer)
			if err == nil && n > 0 {
				incoming <- buffer[:n]
				l.debug("packet received", 2)
			}
		}
	}()
	go func() {
		for true {
			l.mu.RLock()
			delayPackets := l.delayPackets
			lossPackets := l.lossPackets
			duplicatePackets := l.duplicatePackets
			reorderPackets := l.reorderPackets
			corruptPackets := l.corruptPackets
			paused := l.paused
			queued := l.queued
			l.mu.RUnlock()

			if paused {
				continue
			}
			if queued {
				_ = <-l.allowMessageIfQueued
			}

			newPacket := <-incoming
			l.debug("packet dequeued", 2)
			// handle delay
			time.Sleep(time.Duration(delayPackets) * time.Millisecond)

			// handle loss
			if rand.Intn(100) < lossPackets {
				//continue
				l.sendControllerMessage("loss packet is not enabled on in Local_Proxy")
			}

			// handle duplication
			dupC := 1
			if rand.Intn(100) < duplicatePackets {
				//dupC++
				l.sendControllerMessage("Duplicating packet is not enabled in Local_Proxy")
			}

			// handle reordering
			if rand.Intn(100) < reorderPackets {
				//incoming <- newPacket
				//continue
				l.sendControllerMessage("Reordering packet is not enabled in Local_Proxy")
			}

			// handle corruption
			if rand.Intn(100) < corruptPackets {
				l.sendControllerMessage("Corrupting packet not supported in Local_Proxy")
			}

			for i := 0; i < dupC; i++ {
				_, err := dCon.Write(newPacket)
				if err != nil {
					fmt.Println("socket writing error " + dCon.RemoteAddr().String() + " " + err.Error())
					return
				} else {
					l.debug("packet sent", 2)
				}
			}
		}
	}()
}

func (l *Local_Proxy) sendControllerMessage(m string) {
	l.c.SendControllerMessage(&proto.Message{
		StrParams: []string{m},
	})
}

func (l *Local_Proxy) DelayPackets(delay int) error {
	l.mu.Lock()
	l.delayPackets = delay
	l.mu.Unlock()
	l.debug("set new delay", 2)
	return nil

}

func (l *Local_Proxy) LossPackets(loss int) error {
	l.mu.Lock()
	l.lossPackets = loss
	l.mu.Unlock()
	l.debug("set new loss", 2)
	return nil
}

func (l *Local_Proxy) DuplicatePackets(dup int) error {
	l.mu.Lock()
	l.duplicatePackets = dup
	l.mu.Unlock()
	l.debug("set new duplication", 2)
	return nil
}

func (l *Local_Proxy) ReorderPackets(re int) error {
	l.mu.Lock()
	l.reorderPackets = re
	l.mu.Unlock()
	l.debug("set new reorder", 2)
	return nil
}

func (l *Local_Proxy) CorruptPackets(corrupt int) error {
	l.mu.Lock()
	l.corruptPackets = corrupt
	l.mu.Unlock()
	l.debug("set new corrupt", 2)
	return nil
}

func (l *Local_Proxy) Pause(on bool) error {
	l.mu.Lock()
	l.paused = on
	l.mu.Unlock()
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
	l.mu.Lock()
	l.delayPackets = 0
	l.lossPackets = 0
	l.duplicatePackets = 0
	l.reorderPackets = 0
	l.corruptPackets = 0
	l.mu.Unlock()
	l.QueueAllMessages(false)
	l.Pause(false)
	l.debug("reset all", 2)
	return nil
}

func (l *Local_Proxy) Kill() error {
	l.ResetAll()
	l.CleanUp()
	err := util.RunCommand("kill", []string{"-9", l.process_id})
	l.debug("killed", 2)
	return err
}

func (l *Local_Proxy) QueueAllMessages(on bool) error {
	l.mu.Lock()
	l.queued = on
	numThreads := l.numConnections
	l.mu.Unlock()
	for i := 0; i < numThreads+100; i++ {
		if !on {
			l.allowMessageIfQueued <- true
		}
	}
	return nil
}

func (l *Local_Proxy) AllowMessages(n int) error {

	go func() {
		for i := 0; i < n; i++ {
			select {
			case l.allowMessageIfQueued <- true:
				break
			default:
				i--
				break
			}
		}
	}()
	return nil
}

func (l *Local_Proxy) CorruptDB() error {
	l.sendControllerMessage("CorruptDB is not supported by Local_Proxy")
	return nil
}

func (l *Local_Proxy) CleanUp() error {
	for {
		select {
		case l.allowMessageIfQueued <- true:
			break
		default:
			break
		}
	}
}

func (l *Local_Proxy) debug(m string, level int) {
	if l.debugOn && level >= l.debugLevel {
		fmt.Println(m + "\n")
	}
}
