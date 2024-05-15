package torture

import (
	"strconv"
	"toture-test/torture/proto"
)

/*
	the main loop of the controller
*/

func (c *TortureController) Run() {
	go func() {
		for true {
			m_object := <-c.incomingMessages
			c.debug("controller received message from "+strconv.FormatInt(int64(m_object.sender), 10), 0)
			m := (m_object.m).(*proto.Message)
			c.handleMessage(m, m_object.sender)
		}
	}()

}

func (c *TortureController) handleMessage(message *proto.Message, sender int) {

}
