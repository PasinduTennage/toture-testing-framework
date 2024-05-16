package torture

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
	"toture-test/torture/configuration"
	"toture-test/torture/proto"
)

type TortureClient struct {
	name              int           // unique node id
	serverAddress     string        // listening address of self
	controllerAddress string        // ip of the controller
	controllerWriter  *bufio.Writer // IO writer to the controller
	mutex             *sync.Mutex   // IO reader to the controller

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
		mutex:             &sync.Mutex{},
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

// start listening to the client tcp connections

func (c *TortureClient) NetworkInit() {
	c.waitForConnection()
}

/*
	Listen on the client ports for new connection
*/

func (c *TortureClient) waitForConnection() {

	go func(la string) {

		listener, err := net.Listen("tcp", la)
		if err != nil {
			panic(err.Error())
		}
		c.debug("Client listening to messages on "+la, 0)

		conn, err := listener.Accept()
		if err != nil {
			panic(err.Error())
		}
		c.debug("Received incoming tcp connection from controller ", -1)
		go c.controllerListener(bufio.NewReader(conn))
		c.debug("Started listening to controller ", 1)

	}(c.serverAddress)

}

/*
	listen to controller
*/

func (c *TortureClient) controllerListener(reader *bufio.Reader) {

	var err error = nil
	for true {
		obj := (&proto.Message{}).New()
		if err = obj.Unmarshal(reader); err != nil {
			//pr.debug("Error while unmarshalling", 0)
			return
		}
		c.handlerControllerMessage(obj.(*proto.Message))
	}
}

/*
	make TCP connection to controller
*/

func (c *TortureClient) ConnectToController() {

	var b [4]byte
	bs := b[:4]

	for true {
		conn, err := net.Dial("tcp", c.controllerAddress)
		if err == nil {
			c.controllerWriter = bufio.NewWriter(conn)
			c.mutex = &sync.Mutex{}
			binary.LittleEndian.PutUint16(bs, uint16(c.name))
			_, err := conn.Write(bs)
			if err != nil {
				//pr.debug("Error connecting to client "+strconv.Itoa(int(id)), 0)
				panic(err.Error())
			}
			c.debug("Started outgoing tcp connection to controller ", 0)
			break
		} else {
			c.debug("retrying because failed to connect to controller", 0)
			time.Sleep(time.Duration(10) * time.Millisecond)
		}
	}

}

/*
	write a message to a controller
*/

func (c *TortureClient) sendControllerMessage(msg *proto.Message) {

	c.debug("sending message to  controller", 0)

	w := c.controllerWriter
	m := c.mutex

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
	c.debug("sent message to  controller ", 1)
	m.Unlock()
}
