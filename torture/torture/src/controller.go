package torture

import (
	"fmt"
	"time"
	"toture-test/torture/proto"
)

func (c *TortureController) handleMessage(message *proto.Message, sender int) {
	// print message
	fmt.Printf("Controller received message %v from %d\n", message, sender)
}

func (c *TortureController) StartAttack() {
	for true {
		for cl, _ := range c.clients {
			c.sendMessage(int64(cl), &proto.Message{
				Operation: 1,
				IntParams: []int32{1},
				StrParams: []string{"hello"},
			})
		}
		time.Sleep(2 * time.Second)
	}
}
