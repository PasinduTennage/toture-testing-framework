package torture

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
)

type IncomingMessage struct {
	m      proto.Serializable
	sender int
}

type TortureController struct {
	name          int                   // unique node id
	serverAddress string                // listening address of self
	clients       map[int]string        // ip address of each torture client
	clientWriters map[int]*bufio.Writer // IO writer to each client
	mutexes       map[int]*sync.Mutex   // IO writer to each client

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
		clientWriters:    make(map[int]*bufio.Writer),
		mutexes:          make(map[int]*sync.Mutex),
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

	cr.debug("initialed a new torture controller "+strconv.Itoa(cr.name), 0)

	return &cr
}

// debug prints the message to console if the debug is turned on

func (c *TortureController) debug(s string, i int) {
	if c.debugOn && i >= c.debugLevel {
		fmt.Printf("%s\n", s)
	}
}

// start listening to the controller tcp connection

func (c *TortureController) NetworkInit() {
	c.waitForConnections()
}

/*
	Listen on the controller ports for new connections
*/

func (c *TortureController) waitForConnections() {

	go func(la string) {

		var b [4]byte
		bs := b[:4]

		listener, err := net.Listen("tcp", la)
		if err != nil {
			panic(err.Error())
		}
		c.debug("Controller listening to messages on "+la, 0)

		for true {
			conn, err := listener.Accept()
			if err != nil {
				panic(err.Error())
			}
			if _, err := io.ReadFull(conn, bs); err != nil {
				panic(err.Error())
			}
			id := int32(binary.LittleEndian.Uint16(bs))
			c.debug("Received incoming tcp connection from "+strconv.Itoa(int(id)), -1)
			go c.connectionListener(bufio.NewReader(conn), id)
			c.debug("Started listening to "+strconv.Itoa(int(id)), 0)
		}
	}(c.serverAddress)

}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (c *TortureController) connectionListener(reader *bufio.Reader, id int32) {

	var err error = nil
	for true {
		obj := (&proto.Message{}).New()
		if err = obj.Unmarshal(reader); err != nil {
			//pr.debug("Error while unmarshalling", 0)
			return
		}
		c.incomingMessages <- &IncomingMessage{m: obj, sender: int(id)}
	}
}

/*
	make TCP connections to other clients
*/

func (c *TortureController) ConnectToClients() {

	for id, address := range c.clients {

		for true {
			conn, err := net.Dial("tcp", address)
			if err == nil {
				c.clientWriters[id] = bufio.NewWriter(conn)
				c.mutexes[id] = &sync.Mutex{}
				c.debug("Started outgoing tcp connection to client "+address, 0)
				break
			} else {
				c.debug("retrying because failed to connect to "+address, 0)
				time.Sleep(time.Duration(10) * time.Millisecond)
			}
		}

	}

}

/*
	write a message to a client
*/

func (c *TortureController) sendMessage(client int64, msg *proto.Message) {

	c.debug("sending message to  "+strconv.Itoa(int(client)), 0)

	w := c.clientWriters[int(client)]
	m := c.mutexes[int(client)]

	m.Lock()
	err := msg.Marshal(w)
	if err != nil {
		c.debug("Error while marshalling", 0)
		m.Unlock()
		return
	}
	err = w.Flush()
	if err != nil {
		c.debug("Error while flushing", 0)
		m.Unlock()
		return
	}
	c.debug("sent message to  "+strconv.Itoa(int(client)), 1)
	m.Unlock()
}

/*
	the main loop of the controller
*/

func (c *TortureController) Run() {
	go func() {
		for true {
			m_object := <-c.incomingMessages
			c.debug("controller received message from "+strconv.FormatInt(int64(m_object.sender), 10), 0)
			m := (m_object.m).(*proto.Message)
			c.handleMessage(m, m_object.sender)
		}
	}()

}

func (c *TortureController) handleMessage(message *proto.Message, sender int) {

}
