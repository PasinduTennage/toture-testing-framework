package client

import (
	"fmt"
	"sync"
	"time"
	"toture-test/consenbench/common"
	"toture-test/util"
)

type ClientOptions struct {
	NodeInfoFile   string // the yaml file containing the ip address of each node, controller port, client por
	DebugOn        bool
	DebugLevel     int
	LogFileAbsPath string
}

type Client struct {
	Id           int
	Network      *common.Network // to communicate with the controller
	InputChan    chan *common.RPCPairPeer
	logger       *util.Logger
	Options      ClientOptions
	ControllerId int
}

func NewClient(Id int, options ClientOptions) *Client {
	return &Client{
		Id:        Id,
		InputChan: make(chan *common.RPCPairPeer, 10000),
		logger:    util.NewLogger(options.DebugLevel, options.DebugOn, options.LogFileAbsPath),
		Options:   options,
	}
}

// initialize the network layer and run the client

func (c *Client) Run() {

	nodes := common.GetNodes(c.Options.NodeInfoFile)
	listenAddress := ""
	for i := 0; i < len(nodes); i++ {
		if nodes[i].Id == c.Id {
			listenAddress = nodes[i].Ip + ":10080"
		}
	}
	if listenAddress == "" {
		panic("node id not found in the node info file")
	}

	controller := common.GetController(c.Options.NodeInfoFile)
	c.ControllerId = controller.Id

	network_config := common.NetworkConfig{
		ListenAddress:   listenAddress,
		RemoteAddresses: map[int]string{controller.Id: controller.Ip + ":10080"},
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
	c.logger.Debug("Connected to controller, both ways!", 0)

	c.SendStats()

	c.run()

	for true {
		time.Sleep(5 * time.Second)
	}
}

// respond to different messages from the controller

func (c *Client) run() {
	go func() {
		for true {
			rpcPair := <-c.InputChan
			c.logger.Debug(fmt.Sprintf("%v", rpcPair.RpcPair), 0)
			switch rpcPair.RpcPair.Code {
			case common.GetRPCCodes().ControlMsg:
				// handle control message
				ctrlMsg := rpcPair.RpcPair.Obj.(*common.ControlMsg)
				c.Handle(ctrlMsg)
			default:
				c.logger.Debug("Unknown message type", 0)
			}
		}
	}()
}

// periodically send machine stats to the controller

func (c *Client) SendStats() {
	// send machine stats to the controller
	go func() {
		for true {
			// scrape machine stats and send to the controller
			//perf_name := []string{"cpu_usage", "mem_usage", "packetsInRate", "packetsOutRate"}

			cpu := util.GetCPUUsage()
			mem := util.GetMemoryUsage()
			packetsInRate, packetsOutRate := util.GetNetworkStats() // has a 1s sync delay

			perf_stats := []float32{float32(cpu), float32(mem), float32(packetsInRate), float32(packetsOutRate)}
			c.Network.Send(&common.RPCPairPeer{
				RpcPair: &common.RPCPair{
					Code: common.GetRPCCodes().ControlMsg,
					Obj: &common.ControlMsg{
						OperationType: int32(common.GetOperationCodes().Stats),
						FloatArgs:     perf_stats,
					},
				},
				Peer: c.ControllerId,
			})
			time.Sleep(1 * time.Second)
			c.logger.Debug("Sent stats to controller", 0)

		}
	}()
}
