package torture

import (
	"fmt"
	"toture-test/torture/proto"
)

func (cl *TortureClient) handlerControllerMessage(message *proto.Message) {
	fmt.Printf("Client received message from controller %v\n", message)
	cl.sendControllerMessage(message)
}
