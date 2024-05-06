package dummy

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"toture-test/dummy/configuration"
)

type Proxy struct {
	name        int64 // unique node id
	numReplicas int

	addrList        map[int64][]string // map with the IP:port address
	outgoingWriters map[int64][]*bufio.Writer

	serverAddress []string // proxy address

	debugOn    bool // if turned on, the debug messages will be print on the console
	debugLevel int  // debug level

	serverStarted bool // true if the first status message with operation type 1 received

	startTime       time.Time
	receivedLatency []int64

	sent sync.Map

	incomingChan chan *ReceivedMessage

	counter int64
}

type Request struct {
	sentTime time.Time
}

type ReceivedMessage struct {
	message Serializable
	sender  int32
}

func NewProxy(name int64, cfg configuration.InstanceConfig, debugOn bool, debugLevel int) *Proxy {

	pr := Proxy{
		name:            name,
		numReplicas:     len(cfg.Peers),
		addrList:        make(map[int64][]string),
		outgoingWriters: make(map[int64][]*bufio.Writer),
		serverAddress:   []string{},
		debugOn:         debugOn,
		debugLevel:      debugLevel,
		serverStarted:   false,
		incomingChan:    make(chan *ReceivedMessage, 1000),
		counter:         0,
		receivedLatency: make([]int64, 0),
	}

	// initialize the addrList

	for i := 0; i < pr.numReplicas; i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		addresses := make([]string, 0)
		for j := 0; j < len(cfg.Peers[i].PORTS); j++ {
			addresses = append(addresses, cfg.Peers[i].IP+":"+cfg.Peers[i].PORTS[j])
		}
		if pr.name != int64(intName) {
			pr.addrList[int64(intName)] = addresses
			pr.outgoingWriters[int64(intName)] = make([]*bufio.Writer, 0)
		}
		if pr.name == int64(intName) {
			pr.serverAddress = addresses
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())

	pr.debug("initialed a new proxy "+strconv.Itoa(int(pr.name)), 0)

	return &pr
}

/*
	the main loop of the proxy
*/

func (pr *Proxy) Run() {
	go func() {
		for true {
			m_object := <-pr.incomingChan
			pr.debug("Received message", 0)
			pr.handleMessage(m_object.message.(*Message), m_object.sender)
		}
	}()

}

// debug prints the message to console if the debug is turned on

func (pr *Proxy) debug(s string, i int) {
	if pr.debugOn && i >= pr.debugLevel {
		fmt.Printf("%s\n", s)
	}
}
