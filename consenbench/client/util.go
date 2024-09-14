package client

import (
	"fmt"
	"time"
	"toture-test/consenbench/common"
	"toture-test/util"
)

// Set the replica name and the ports that the replica is listening on

func (c *Client) SetPorts(msg *common.ControlMsg) {
	replica_name := msg.StringArgs[0]
	ports := msg.StringArgs[1:]
	c.Attacker.Process_name = replica_name
	c.Attacker.Ports_under_attack = ports
	c.logger.Debug("Set replica name to "+replica_name+" and ports to "+fmt.Sprintf("%v", ports), 3)
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
			time.Sleep(500 * time.Millisecond)
			c.logger.Debug("Sent stats to controller", 0)

		}
	}()
}
