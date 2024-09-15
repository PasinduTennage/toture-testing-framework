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
	c.logger.Debug("Initialized TC ", 3)
}

// Execute the Pending commands

func (c *Client) ExecuteLastNetEmCommands() error {
	var err error
	for i := 0; i < len(c.Attacker.NextNetEmCommands); i++ {
		if len(c.Attacker.NextNetEmCommands[i]) == 0 {
			continue
		}
		err = c.RunCommand(c.Attacker.NextNetEmCommands[i][0], c.Attacker.NextNetEmCommands[i][1:])
	}
	c.Attacker.NextNetEmCommands = [][]string{}
	return err
}

// apply the handle to each port

func (c *Client) applyHandleToEachPort() {
	i := 1
	for _, port := range c.Attacker.Ports_under_attack {
		c.RunCommand("tc", []string{"filter", "add", "dev", c.Attacker.Device, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(i), "u32", "match", "ip", "dport", port, "0xffff", "flowid", c.Attacker.Parent_band})
		c.Attacker.NextNetEmCommands = append(c.Attacker.NextNetEmCommands, []string{"tc", "filter", "del", "dev", c.Attacker.Device, "protocol", "ip", "parent", "1:0", "prio", strconv.Itoa(i), "u32", "match", "ip", "dport", port, "0xffff", "flowid", c.Attacker.Parent_band})
		i++
	}
}

func (c *Client) decrementDelay() {
	c.Attacker.DelayPackets--
}

func (c *Client) SetNewHandler() error {
	if c.Attacker.ReorderPackets > 0 && c.Attacker.DelayPackets == 0 {
		c.Attacker.DelayPackets = 1
		defer c.decrementDelay()
	}

	err := c.RunCommand("tc", []string{"qdisc", "add", "dev", c.Attacker.Device, "parent", c.Attacker.Parent_band, "handle", c.Attacker.Handle + ":", "netem", "delay", strconv.Itoa(c.Attacker.DelayPackets) + "ms", "loss", strconv.Itoa(c.Attacker.LossPackets) + "%", "duplicate", strconv.Itoa(c.Attacker.DuplicatePackets) + "%", "reorder", strconv.Itoa(c.Attacker.ReorderPackets) + "%", "50%", "corrupt", strconv.Itoa(c.Attacker.CorruptPackets) + "%"})
	c.applyHandleToEachPort()
	c.Attacker.NextNetEmCommands = append(c.Attacker.NextNetEmCommands, []string{"tc", "qdisc", "del", "dev", c.Attacker.Device, "parent", c.Attacker.Parent_band, "handle", c.Attacker.Handle + ":", "netem", "delay", strconv.Itoa(c.Attacker.DelayPackets) + "ms", "loss", strconv.Itoa(c.Attacker.LossPackets) + "%", "duplicate", strconv.Itoa(c.Attacker.DuplicatePackets) + "%", "reorder", strconv.Itoa(c.Attacker.ReorderPackets) + "%", "50%", "corrupt", strconv.Itoa(c.Attacker.CorruptPackets) + "%"})
	c.logger.Debug("Set new net em handler", 3)
	return err
}

func (c *Client) CleanUp() error {
	c.RunCommand("tc", []string{"filter", "del", "dev", c.Attacker.Device})
	c.RunCommand("tc", []string{"qdisc", "del", "dev", c.Attacker.Device, "root"})
	return nil
}

// kill the client

func (c *Client) Kill() {
	c.ExecuteLastNetEmCommands()
	c.CleanUp()
	c.RunCommand("pkill", []string{c.Attacker.Process_name})
	c.logger.Debug("killed consensus node", 3)
}

func (c *Client) intern_slowdown() {
	slowdown := false
	for true {
		select {
		case on := <-c.Attacker.On_Off_Chan:
			if on {
				slowdown = true
			} else {
				slowdown = false
			}
		default:
			slowdown = slowdown

		}
		if slowdown {
			c.Pause()
			time.Sleep(100 * time.Millisecond)
			c.Continue()
			c.logger.Debug("slowdowned", 3)
		}
	}
}

// slow down the client

func (c *Client) SlowDown(action string) {
	if action == "true" {
		c.Attacker.On_Off_Chan <- true
	} else {
		c.Attacker.On_Off_Chan <- false
	}
}

// pause the client

func (c *Client) Pause() {
	c.RunCommand("pkill", []string{"-STOP", c.Attacker.Process_name})
	c.logger.Debug("paused", 2)
}

// continue the client

func (c *Client) Continue() {
	c.RunCommand("pkill", []string{"-CONT", c.Attacker.Process_name})
	c.logger.Debug("continue", 2)
}

// set the skew

func (c *Client) SetSkew(f float32) {
	// TODO
	panic("Not implemented")
}

// set the drift

func (c *Client) SetDrift(f float32) {
	panic("Not implemented")
}

// set the delay

func (c *Client) SetDelay(f float32) {
	c.ExecuteLastNetEmCommands()
	c.Attacker.DelayPackets = int(f)
	c.logger.Debug("set delay", 3)
	c.SetNewHandler()
}

// set the loss

func (c *Client) SetLoss(f float32) {
	c.ExecuteLastNetEmCommands()
	c.Attacker.LossPackets = int(f)
	c.logger.Debug("set loss", 3)
	c.SetNewHandler()
}

// set the bandwidth

func (c *Client) SetBandwidth(f float32) {
	// TODO
	panic("Not implemented")
}
