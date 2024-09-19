package client

import "time"

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
			c.logger.Debug("Slowdowned inside thread", 3)
		}
	}
}

// slow down the client

func (c *Client) SlowDown(action string) {
	if action == "true" {
		select {
		case c.Attacker.On_Off_Chan <- true:
			c.logger.Debug("slowdown", 3)
		default:
			c.logger.Debug("cannot invoke slowdown -- buffers filled", 3)
		}
	} else {
		select {
		case c.Attacker.On_Off_Chan <- false:
			c.logger.Debug("Cancelled slowdown", 3)
		default:
			c.logger.Debug("cannot cancel slowdown -- buffers filled", 3)
		}
	}
}

// pause the client

func (c *Client) Pause() {
	RunCommand("pkill", []string{"-STOP", "-f", c.Attacker.Process_name}, c.logger)
	c.logger.Debug("paused", 3)
}

// continue the client

func (c *Client) Continue() {
	RunCommand("pkill", []string{"-CONT", "-f", c.Attacker.Process_name}, c.logger)
	c.logger.Debug("continue", 3)
}

// kill the client

func (c *Client) Kill() {
	for _, v := range c.Attacker.NetEmAttackers {
		v.ExecuteLastNetEmCommands()
	}
	c.CleanUp()
	RunCommand("pkill", []string{"-KILL", "-f", c.Attacker.Process_name}, c.logger)
	c.logger.Debug("killed consensus node", 3)
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
