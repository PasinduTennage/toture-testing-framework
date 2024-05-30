package torture

import (
	"fmt"
	"os"
	"toture-test/torture/proto"
)

func (cl *TortureClient) handlerControllerMessage(message *proto.Message) {
	cl.debug(fmt.Sprintf("Client received message from controller %v\n", message), 0)

	if message.Operation == NewOperationTypes().DelayPackets {
		cl.attacker.DelayPackets(int(message.IntParams[0]), message.On)
	}

	if message.Operation == NewOperationTypes().LossPackets {
		cl.attacker.LossPackets(int(message.IntParams[0]), message.On)
	}

	if message.Operation == NewOperationTypes().DuplicatePackets {
		cl.attacker.DuplicatePackets(int(message.IntParams[0]), message.On)
	}

	if message.Operation == NewOperationTypes().ReorderPackets {
		cl.attacker.ReorderPackets(int(message.IntParams[0]), message.On)
	}

	if message.Operation == NewOperationTypes().CorruptPackets {
		cl.attacker.CorruptPackets(int(message.IntParams[0]), message.On)
	}

	if message.Operation == NewOperationTypes().Pause {
		cl.attacker.Pause(message.On)
	}

	if message.Operation == NewOperationTypes().ResetAll {
		cl.attacker.ResetAll()
	}

	if message.Operation == NewOperationTypes().Kill {
		cl.attacker.Kill()
	}

	if message.Operation == NewOperationTypes().QueueAllMessages {
		cl.attacker.QueueAllMessages(message.On)
	}

	if message.Operation == NewOperationTypes().AllowMessages {
		cl.attacker.AllowMessages(int(message.IntParams[0]))
	}

	if message.Operation == NewOperationTypes().CorruptDB {
		cl.attacker.CorruptDB()
	}

	if message.Operation == EXIT {
		cl.attacker.ResetAll()
		os.Exit(0)
	}

}
