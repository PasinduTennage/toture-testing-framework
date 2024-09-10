package controller

import (
	"fmt"
	"time"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type ControllerOptions struct {
	AttackDuration int      // second
	Attacks        []string // set of attacks to run
	NodeInfoFile   string   // the yaml file containing the ip address of each node, controller port, client port
	debugOn        bool
	debugLevel     int
}

// Controller struct
type Controller struct {
	Id        int
	Nodes     []*common.Node
	Network   *common.Network         // to communicate with the clients
	InputChan chan common.RPCPairPeer // message from the clients
	Options   ControllerOptions
	logger    *util.Logger
}

func NewController(Id int, Options ControllerOptions) *Controller {
	return &Controller{
		Id:        Id,
		InputChan: make(chan common.RPCPairPeer, 10000),
		Options:   Options,
		logger:    util.NewLogger(Options.debugLevel, Options.debugOn),
	}
}

// copy clients and start the client binary, check connections and close clients

func (c *Controller) BootstrapClients() error {

	c.InitiliazeNodes()

	// copy the client binary to all the nodes
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].ExecCmd(fmt.Sprintf("mkdir -p %vbench", c.Nodes[i].HomeDir))
		c.Nodes[i].Put_Load("consenbench/bin/bench", fmt.Sprintf("%vbench/", c.Nodes[i].HomeDir))
		c.Nodes[i].Put_Load("consenbench/assets/ip.yaml", fmt.Sprintf("%vbench/", c.Nodes[i].HomeDir))
	}

	c.logger.Debug("Copied the client binary to all the nodes", 0)

	// start the client binary
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].Start_Client()
	}
	time.Sleep(5 * time.Second)
	c.logger.Debug("Started the client binary on all the nodes", 0)

	// initiate the tcp connections
	c.NetworkInit()
	c.logger.Debug("Initialized the network layer with all clients", 0)

	// close the clients
	defer c.CloseClients()
	c.logger.Debug("Closed the clients", 0)

	c.logger.Debug("Bootstrapped the clients, exiting", 0)
	return nil

}

// start listening to incoming connections and connect to all remote nodes

func (c *Controller) NetworkInit() error {
	// initialize the network layer
	return nil
}

// copy the consensus binary

func (c *Controller) CopyConsensus(protocol string) error {
	// copy the consensus binary to all the nodes
	return nil
}

// run the controller

func (c *Controller) Run(protocol string) error {

	// start all remote clients using node interface and tcp connect

	for i := 0; i < len(c.Options.Attacks); i++ {
		// bootstrap the consensus protocol

		// instantiate the node and link objects and start the attack, only instantiate for the attack_nodes that are running replicas, use the consensus options object

		// collect stats from consensus object and print them
	}

	return nil
}

func (c *Controller) HandleInputStream() error {
	// handle the messages from the clients about machine stats
	go func() {
		for true {
			_ = <-c.InputChan
			// update the stats of the node
		}
	}()
	return nil
}

func (c *Controller) InitiliazeNodes() {

}

func (c *Controller) CloseClients() {
	c.Network.Broadcast(&common.RPCPair{
		Code: common.GetRPCCodes().ControlMsg,
		Obj: &common.ControlMsg{
			OperationType: int32(common.GetOperationCodes().ShutDown),
			StringArgs:    nil,
			IntArgs:       nil,
		},
	})
}
