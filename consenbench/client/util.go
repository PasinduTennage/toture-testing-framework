package client

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"toture-test/consenbench/common"
	"toture-test/util"
)

// Set the replica name and the ports that the replica is listening on

func (c *Client) SetPorts(msg *common.ControlMsg) {
	replica_name := msg.StringArgs[0]
	ports := msg.StringArgs[1:]
	if len(ports) == 0 {
		panic("No ports provided")
	}
	c.Attacker.Process_name = replica_name
	c.Attacker.Ports_under_attack = ports
	c.logger.Debug("Set replica name to "+replica_name+" and ports to "+fmt.Sprintf("%v", ports), 3)
	c.Init()
}

// periodically send machine stats to the controller

func (c *Client) SendStats() {
	// send machine stats to the controller
	go func() {
		for true {
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
			time.Sleep(100 * time.Millisecond)
			c.logger.Debug("Sent stats to controller", 0)

		}
	}()
}

// runCommand runs the given command with the provided arguments

func (c *Client) RunCommand(name string, arg []string) error {
	cmd := exec.Command(name, arg...)
	if cmd.Err != nil {
		fmt.Println("Error running command " + name + " " + strings.Join(arg, " ") + " " + cmd.Err.Error() + "\n")
		return cmd.Err
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running command " + name + " " + strings.Join(arg, " ") + " " + err.Error() + "\n")
		return err
	} else {
		c.logger.Debug("Success command "+name+" "+strings.Join(arg, " "), 3)
	}
	return nil
}

// Initialize the netem client

func (c *Client) Init() {
	c.RunCommand("tc", []string{"filter", "del", "dev", c.Options.Device})
	c.RunCommand("tc", []string{"qdisc", "del", "dev", c.Options.Device, "root"})
	c.RunCommand("tc", []string{"qdisc", "add", "dev", c.Options.Device, "root", "handle", "1:", "prio", "bands", strconv.Itoa(5)})
	go c.intern_slowdown()
	c.logger.Debug("Initialized TC ", 3)
}
