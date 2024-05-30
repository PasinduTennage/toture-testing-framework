package torture

/*
	Steps to add a new attack
	1. Add the attack method in the Attacker interface in attacker_interface.go
	2. Add a new field to OperationTypes and assign an int in NewOperationTypes in the attacker_interface.go
	3. Add a new case in the handlerControllerMessage method in client.go
	4. Add a new method in controller_node.go







*/

// Attacker interface defines the methods that an attacker should implement
type Attacker interface {
	DelayPackets(int, bool) error     // delay messages by delay milliseconds
	LossPackets(int, bool) error      // loss percentage of packets
	DuplicatePackets(int, bool) error // duplicate percentage of packets
	ReorderPackets(int, bool) error   // reorder percentage of packets
	CorruptPackets(int, bool) error   // corrupt percentage of packets
	Pause(bool) error                 // halt the process
	ResetAll() error                  // reset the process
	Kill() error                      // kill the process
	QueueAllMessages(bool) error      // buffer all messages
	AllowMessages(int) error          // allow num_messages messages
	CorruptDB() error                 // corrupt the internal database
}

type OperationTypes struct {
	DelayPackets     int32
	LossPackets      int32
	DuplicatePackets int32
	ReorderPackets   int32
	CorruptPackets   int32
	Pause            int32
	ResetAll         int32
	Kill             int32
	QueueAllMessages int32
	AllowMessages    int32
	CorruptDB        int32
}

func NewOperationTypes() OperationTypes {
	return OperationTypes{
		DelayPackets:     1,
		LossPackets:      2,
		DuplicatePackets: 3,
		ReorderPackets:   4,
		CorruptPackets:   5,
		Pause:            6,
		ResetAll:         7,
		Kill:             8,
		QueueAllMessages: 9,
		AllowMessages:    10,
		CorruptDB:        11,
	}
}

const EXIT = 255
