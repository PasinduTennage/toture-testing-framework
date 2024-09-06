package controller

import (
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type ControllerOptions struct {
	AttackDuration int      // second
	Attack         []string // set of attacks to run
	NodeInfoFile   string   // the yaml file containing the ip address of each node, controller port, client port
}

// Controller struct
type Controller struct {
	Id         int
	Nodes      []*util.Node
	Network    *common.Network // to communicate with the clients
	Consensus  protocols.Consensus
	InputChan  chan common.OutgoingRPC
	OutputChan chan common.OutgoingRPC
}

func NewController(options ControllerOptions) *Controller {
	return &Controller{}
}

func (c *Controller) NetworkInit() error {
	// initialize the network layer
	return nil
}

// copy clients and start the client binary

func (c *Controller) BootstrapClients() error {
	// copy the client binary to all the nodes
	// start the client binary and initiate the tcp connections
	// if all tcp connections succeed than end the remote client program
	return nil
}

// copy the consensus binary

func (c *Controller) CopyConsensus(protocol string) error {
	// copy the consensus binary to all the nodes
	return nil
}

// run the controller

func (c *Controller) Run(protocol string, attack string) error {
	// initialize the consensus object using the protocol string

	// run all remote clients using node interface and tcp connect

	// bootstrap the consensus protocol

	// instantiate the node and link objects and start the attack, only instantiate for the nodes that are running replicas, use the consensus options object

	// collect stats from consensus object and print them
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
