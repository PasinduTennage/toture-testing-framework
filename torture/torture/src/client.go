package torture

import (
	"fmt"
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

	if message.Operation == NewOperationTypes().ReorderPercentagePackets {
		cl.attacker.ReorderPercentagePackets(int(message.IntParams[0]))
	}

	if message.Operation == NewOperationTypes().CorruptPercentagePackets {
		cl.attacker.CorruptPercentagePackets(int(message.IntParams[0]))
	}

	if message.Operation == NewOperationTypes().Halt {
		cl.attacker.Halt()
	}

	if message.Operation == NewOperationTypes().Reset {
		cl.attacker.Reset()
	}

	if message.Operation == NewOperationTypes().Kill {
		cl.attacker.Kill()
	}

	if message.Operation == NewOperationTypes().BufferAllMessages {
		cl.attacker.BufferAllMessages()
	}

	if message.Operation == NewOperationTypes().AllowMessages {
		cl.attacker.AllowMessages(int(message.IntParams[0]))
	}

	if message.Operation == NewOperationTypes().CorruptDB {
		cl.attacker.CorruptDB()
	}

	if message.Operation == NewOperationTypes().Exit {
		cl.attacker.Exit()
	}

}
