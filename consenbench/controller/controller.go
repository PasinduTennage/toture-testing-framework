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
	InputChan  chan interface{}
	OutputChan chan interface{}
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

func (c *Controller) Run(protocol string) error {
	// initialize the consensus object using the protocol string

	// run all remote clients using node interface and tcp connect

	// bootstrap the consensus protocol
	//c.consensus.Bootstrap(c.nodes, c.consensus.ExtractOptions())
	//
	return nil
}
