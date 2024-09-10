package client

import (
	"os"
	"toture-test/consenbench/common"
)

func (c *Client) Handle(msg *common.ControlMsg) {
	if int(msg.OperationType) == common.GetOperationCodes().ShutDown {
		c.logger.Debug("Received exit signal from controller", 0)
		os.Exit(0)
	} else if int(msg.OperationType) == common.GetOperationCodes().Stats {
		panic("Received stats from controller, should not happen")
	}
	// todo handle other messages
}
