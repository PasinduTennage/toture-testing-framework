package client

import "strconv"

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
