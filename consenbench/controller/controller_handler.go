package controller

import (
	"fmt"
	"toture-test/consenbench/common"
)

func (c *Controller) Handle(msg *common.ControlMsg, node int) {
	if msg.OperationType == int32(common.GetOperationCodes().Stats) {
		c.logger.Debug(fmt.Sprintf("Received stats %v from %v", msg.FloatArgs, node), 0)
		c.Nodes[node-1].UpdateStats(msg.FloatArgs)
	} else {
		c.logger.Debug(fmt.Sprintf("Unknown operation type %v", msg.OperationType), 0)
	}
}
