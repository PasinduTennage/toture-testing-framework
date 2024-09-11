package controller

import (
	"fmt"
	"os"
	"sync"
	"time"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type ControllerOptions struct {
	AttackDuration int      // second
	Attacks        []string // set of attacks to run
	NodeInfoFile   string   // the yaml file containing the ip address of each node, controller port, client port
	DebugOn        bool
	DebugLevel     int
}

// Controller struct
type Controller struct {
	Id        int
	Nodes     []*common.Node
	Network   *common.Network          // to communicate with the clients
	InputChan chan *common.RPCPairPeer // message from the clients
	Options   ControllerOptions
	logger    *util.Logger
}

func NewController(Id int, Options ControllerOptions) *Controller {
	return &Controller{
		Id:        Id,
		InputChan: make(chan *common.RPCPairPeer, 10000),
		Options:   Options,
		logger:    util.NewLogger(Options.DebugLevel, Options.DebugOn),
	}
}

// copy clients and start the client binary, check connections and close clients

func (c *Controller) BootstrapClients() error {

	c.InitiliazeNodes()

	// copy the client binary to all the nodes
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].ExecCmd(fmt.Sprintf("pkill -f bench && rm -r %vbench && mkdir -p %vbench", c.Nodes[i].HomeDir, c.Nodes[i].HomeDir))
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

	c.HandleClientMessages()

	// close the clients
	defer c.CloseClients()
	c.logger.Debug("Closed the clients", 0)

	c.logger.Debug("Bootstrapped the clients, exiting", 0)
	os.Exit(0)
	return nil

}

// start listening to incoming connections and connect to all remote nodes

func (c *Controller) NetworkInit() error {
	network_config := common.NetworkConfig{
		ListenAddress:   common.GetController(c.Options.NodeInfoFile).Ip + ":10080",
		RemoteAddresses: common.GetRemoteAddresses(c.Nodes),
	}
	c.Network = common.NewNetwork(c.Id, &network_config, c.InputChan, c.logger)
	c.Network.RegisterRPC(&common.ControlMsg{}, common.GetRPCCodes().ControlMsg)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		c.Network.ConnectRemotes()
		wg.Done()
	}()
	c.Network.Listen()
	wg.Wait()
	c.logger.Debug("Connected to all remote nodes, both ways!", 0)
	return nil
}

// copy the consensus binary

func (c *Controller) CopyConsensus(protocol string) {
	// copy the consensus binary to all the nodes
}

// run the controller

func (c *Controller) Run(protocol string) {

	// start all remote clients using node interface and tcp connect

	for i := 0; i < len(c.Options.Attacks); i++ {
		// bootstrap the consensus protocol

		// instantiate the node and link objects and start the attack, only instantiate for the attack_nodes that are running replicas, use the consensus options object

		// collect stats from consensus object and print them
	}

}

func (c *Controller) HandleClientMessages() error {
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
	c.Nodes = common.GetNodes(c.Options.NodeInfoFile)
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].InitNode(c.logger)
	}
}

func (c *Controller) CloseClients() {
	c.Network.Broadcast(&common.RPCPair{
		Code: common.GetRPCCodes().ControlMsg,
		Obj: &common.ControlMsg{
			OperationType: int32(common.GetOperationCodes().ShutDown),
		},
	})
}
