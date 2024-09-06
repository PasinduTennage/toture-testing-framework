package controller

import (
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type ControllerOptions struct {
	AttackDuration int      // second
	Attacks        []string // set of attacks to run
	NodeInfoFile   string   // the yaml file containing the ip address of each node, controller port, client port
}

// Controller struct
type Controller struct {
	Id         int
	Nodes      []*util.Node
	Network    *common.Network // to communicate with the clients
	Consensus  protocols.Consensus
	InputChan  chan common.RPCPairPeer
	OutputChan chan common.RPCPairPeer
	Options    ControllerOptions
}

func NewController(options ControllerOptions) *Controller {
	return &Controller{
		Options: options,
	}
}

// copy clients and start the client binary, check connections and close clients

func (c *Controller) BootstrapClients() error {
	// copy the client binary to all the nodes
	// start the client binary and initiate the tcp connections
	// if all tcp connections succeed than end the remote client program
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
