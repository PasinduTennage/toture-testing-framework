package dummy

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
	"toture-test/dummy/configuration"
)

type Proxy struct {
	name        int64 // unique node identifier as defined in the configuration.yml
	numReplicas int

	addrList        map[int64][]string // map with the IP:port address of every client
	incomingReaders map[int64][]*bufio.Reader
	outgoingWriters map[int64][]*bufio.Writer

	serverAddress []string       // proxy address
	Listeners     []net.Listener // tcp listener for clients

	debugOn    bool // if turned on, the debug messages will be print on the console
	debugLevel int  // debug level

	serverStarted bool // true if the first status message with operation type 1 received

	startTime     time.Time
	receivedCount int
	sentCount     int

	incomingChan chan Serializable
	outgoingChan chan Serializable
}

func NewProxy(name int64, cfg configuration.InstanceConfig, debugOn bool, debugLevel int) *Proxy {

	pr := Proxy{
		name:            name,
		numReplicas:     len(cfg.Peers),
		addrList:        make(map[int64][]string),
		incomingReaders: make(map[int64][]*bufio.Reader),
		outgoingWriters: make(map[int64][]*bufio.Writer),
		serverAddress:   []string{},
		Listeners:       []net.Listener{},
		debugOn:         debugOn,
		debugLevel:      debugLevel,
		serverStarted:   false,
		receivedCount:   0,
		sentCount:       0,
		incomingChan:    make(chan Serializable, 1000),
		outgoingChan:    make(chan Serializable, 1000),
	}

	// initialize the clientAddrList

	for i := 0; i < pr.numReplicas; i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		addresses := make([]string, 0)
		for j := 0; j < len(cfg.Peers[i].PORTS); j++ {
			addresses = append(addresses, cfg.Peers[i].IP+":"+cfg.Peers[i].PORTS[j])
		}
		pr.addrList[int64(intName)] = addresses
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
			m := <-pr.incomingChan
			pr.debug("Received message", 0)
			message := m.(*Message)
			pr.handleMessage(message)
		}
	}()

}

// debug prints the message to console if the debug is turned on

func (pr *Proxy) debug(s string, i int) {
	if pr.debugOn && i >= pr.debugLevel {
		fmt.Printf("%s\n", s)
	}
}

func (pr *Proxy) handleMessage(message *Message) {

}
