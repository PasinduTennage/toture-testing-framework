package controller

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"toture-test/consenbench/common"
	"toture-test/protocols"
	consensus "toture-test/protocols/baxos"
	"toture-test/util"
)

type ControllerOptions struct {
	AttackDuration int    // second
	Attack         string // set of attacks to run
	NodeInfoFile   string // the yaml file containing the ip address of each node, controller port, client port
	DebugOn        bool
	DebugLevel     int
	LogFilePath    string
	Device         string
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
		logger:    util.NewLogger(Options.DebugLevel, Options.DebugOn, Options.LogFilePath),
	}
}

// copy clients and start the client binary, check connections and close clients

func (c *Controller) BootstrapClients() error {

	c.InitiliazeNodes()

	// copy the client binary to all the nodes
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].ExecCmd(fmt.Sprintf("sudo apt update"))
		c.Nodes[i].ExecCmd(fmt.Sprintf("sudo apt install iproute2"))
		c.Nodes[i].ExecCmd(fmt.Sprintf("sudo setcap cap_net_admin,cap_net_raw+ep $(which tc)"))
		c.Nodes[i].ExecCmd(fmt.Sprintf("getcap $(which tc)"))

		c.Nodes[i].ExecCmd(fmt.Sprintf("pkill  bench"))
		c.Nodes[i].ExecCmd(fmt.Sprintf("rm -r %vbench", c.Nodes[i].HomeDir))
		c.Nodes[i].ExecCmd(fmt.Sprintf("mkdir -p %vbench", c.Nodes[i].HomeDir))
		c.Nodes[i].Put_Load("consenbench/bin/bench", fmt.Sprintf("%vbench/", c.Nodes[i].HomeDir))
		c.Nodes[i].Put_Load("consenbench/assets/ip.yaml", fmt.Sprintf("%vbench/", c.Nodes[i].HomeDir))
	}

	fmt.Println("Copied the client binary to all the nodes")

	// start the client binary
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].Start_Client(c.Options.Device)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Started the client binary on all the nodes")

	// initiate the tcp connections
	c.NetworkInit()
	fmt.Println("Initialized the network layer with all clients")

	c.HandleClientMessages()

	time.Sleep(10 * time.Second)
	// close the clients
	c.CloseClients()
	fmt.Println("Closed the clients")
	c.DownloadClientLogs()
	fmt.Println("Downloaded the logs from the clients")
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
	c.InitiliazeNodes()
	var protocol_impl protocols.Consensus
	if protocol == "baxos" {
		protocol_impl = consensus.NewBaxos(c.logger)
	} else {
		panic("Unknown protocol")
	}
	protocol_impl.ExtractOptions("protocols/" + protocol + "/assets/options.yaml")
	protocol_impl.CopyConsensus(c.Nodes)

}

// run the controller

func (c *Controller) Run(protocol string) {
	c.InitiliazeNodes()
	// start the client binary
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].Start_Client(c.Options.Device)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("Started the client binary on all the nodes")

	// initiate the tcp connections
	c.NetworkInit()
	fmt.Println("Initialized the network layer with all clients")

	c.HandleClientMessages()
	time.Sleep(10 * time.Second)

	var protocol_impl protocols.Consensus
	if protocol == "baxos" {
		protocol_impl = consensus.NewBaxos(c.logger)
	} else {
		panic("Unknown protocol")
	}
	protocol_impl.ExtractOptions("protocols/" + protocol + "/assets/options.yaml")

	bootstrap_complete_chan := make(chan bool)
	performance_output_chan := make(chan util.Performance)
	num_replicas_chan := make(chan int)
	process_name_chan := make(chan string)

	go protocol_impl.Bootstrap(c.Nodes, c.Options.AttackDuration, performance_output_chan, bootstrap_complete_chan, num_replicas_chan, process_name_chan)

	num_replicas := <-num_replicas_chan
	process_name_and_ports := <-process_name_chan

	process_name := strings.Split(process_name_and_ports, ",")[0]
	ports := strings.Split(process_name_and_ports, ",")[1:]

	attackNodes, attackLinks, leaderOracle := GetAttackObjects(num_replicas, process_name, c.Nodes, c, c.logger, ports)
	var attack_impl Attack
	if c.Options.Attack == "basic" {
		attack_impl = NewBasicAttack(c.logger)
	} else {
		panic("Unknown attack")

	}

	<-bootstrap_complete_chan // wait for the bootstrap to complete
	fmt.Print("Bootstrap complete, starting attack from controller\n")

	attack_impl.Attack(attackNodes, attackLinks, leaderOracle, c.Options.AttackDuration)
	<-performance_output_chan
	fmt.Print("Attack complete\n")
	c.CloseClients()
	fmt.Println("Closed the clients")
	c.DownloadClientLogs()
	fmt.Println("Downloaded the logs from the clients")
	fmt.Println("test complete")
}

func (c *Controller) HandleClientMessages() error {
	// handle the messages from the clients about machine stats
	go func() {
		for true {
			msg := <-c.InputChan
			c.logger.Debug(fmt.Sprintf("received from client %v %v ", msg.Peer, msg.RpcPair), 0)
			switch msg.RpcPair.Code {
			case common.GetRPCCodes().ControlMsg:
				// handle control message
				ctrlMsg := msg.RpcPair.Obj.(*common.ControlMsg)
				c.Handle(ctrlMsg, msg.Peer)
			default:
				c.logger.Debug("Unknown message type", 0)
			}

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

func (c *Controller) DownloadClientLogs() {
	for i := 0; i < len(c.Nodes); i++ {
		c.Nodes[i].Get_Load(c.Nodes[i].HomeDir+"bench/log.log", fmt.Sprintf("bench/%v.log", c.Nodes[i].Id))
	}
}
