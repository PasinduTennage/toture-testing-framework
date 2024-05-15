package torture

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
)

type TortureClient struct {
	name              int           // unique node id
	serverAddress     string        // listening address of self
	controllerAddress string        // ip of the controller
	controllerWriter  *bufio.Writer // IO writer to the controller
	controllerReader  *bufio.Reader // IO reader to the controller

	debugOn    bool // if turned on, the debug messages will be print on the console
	debugLevel int  // debug level
}

// NewCLient creates a new torture client

func NewClient(name int, cfg configuration.InstanceConfig, debugOn bool, debugLevel int) *TortureClient {

	cl := TortureClient{
		name:              name,
		serverAddress:     "",
		controllerAddress: cfg.Controller.IP + ":" + cfg.Controller.PORT,
		controllerWriter:  nil,
		controllerReader:  nil,
		debugOn:           debugOn,
		debugLevel:        debugLevel,
	}

	// initialize the server address

	for i := 0; i < len(cfg.Peers); i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		if cl.name == intName {
			cl.serverAddress = cfg.Peers[i].IP + ":" + cfg.Peers[i].PORT
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())

	cl.debug("initialed a new torture client "+strconv.Itoa(int(cl.name)), 0)

	return &cl
}

// debug prints the message to console if the debug is turned on

func (cl *TortureClient) debug(s string, i int) {
	if cl.debugOn && i >= cl.debugLevel {
		fmt.Printf("%s\n", s)
	}
}

type IncomingMessage struct {
	m      proto.Serializable
	sender int
}

type TortureController struct {
	name          int                   // unique node id
	serverAddress string                // listening address of self
	clients       map[int]string        // ip address of each torture client
	clientWriters map[int]*bufio.Writer // IO writer to each client
	clientReaders map[int]*bufio.Reader // IO readers for each client

	debugOn    bool // if turned on, the debug messages will be print on the console
	debugLevel int  // debug level

	incomingMessages chan *IncomingMessage
}

// NewController creates a new torture controller

func NewController(name int, cfg configuration.InstanceConfig, debugOn bool, debugLevel int) *TortureController {

	cr := TortureController{
		name:             name,
		serverAddress:    cfg.Controller.IP + ":" + cfg.Controller.PORT,
		clients:          make(map[int]string),
		clientReaders:    make(map[int]*bufio.Reader),
		clientWriters:    make(map[int]*bufio.Writer),
		debugOn:          debugOn,
		debugLevel:       debugLevel,
		incomingMessages: make(chan *IncomingMessage, 10000),
	}

	// initialize the client addresses

	for i := 0; i < len(cfg.Peers); i++ {
		intName, _ := strconv.Atoi(cfg.Peers[i].Name)
		cr.clients[intName] = cfg.Peers[i].IP + ":" + cfg.Peers[i].PORT
	}

	rand.Seed(time.Now().UTC().UnixNano())

	cr.debug("initialed a new torture controller "+strconv.Itoa(int(cr.name)), 0)

	return &cr
}

// debug prints the message to console if the debug is turned on

func (c *TortureController) debug(s string, i int) {
	if c.debugOn && i >= c.debugLevel {
		fmt.Printf("%s\n", s)
	}
}
