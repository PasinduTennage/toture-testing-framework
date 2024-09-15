package controller

import (
	"fmt"
	"toture-test/consenbench/common"
)

func (c *Controller) Handle(msg *common.ControlMsg, node int) {
	if msg.OperationType == int32(common.GetOperationCodes().Stats) {
		c.logger.Debug(fmt.Sprintf("Received stats %v from %v", msg.FloatArgs, node), 3)
		c.Nodes[node-2].UpdateStats(msg.FloatArgs)
	}
}
