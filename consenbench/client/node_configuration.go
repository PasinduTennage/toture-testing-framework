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
	c.RunCommand("pkill", []string{"-STOP", c.Attacker.Process_name})
	c.logger.Debug("paused", 3)
}

// continue the client

func (c *Client) Continue() {
	c.RunCommand("pkill", []string{"-CONT", c.Attacker.Process_name})
	c.logger.Debug("continue", 3)
}

// kill the client

func (c *Client) Kill() {
	c.ExecuteLastNetEmCommands()
	c.CleanUp()
	c.RunCommand("pkill", []string{c.Attacker.Process_name})
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
