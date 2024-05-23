package torture

/*
	Steps to add a new attack
	1. Add the attack method in the Attacker interface in client_attacker_interface.go
	2. Add a new field to OperationTypes and assign an int in NewOperationTypes in the client_attacker_interface.go
	3. In each client attacker implementation (for example netEmAttacker.go), implement the attack method
	4. In client.go add a new if block in handlerControllerMessage to call the attack method
	5. Add a new method in controller_node.go
	6. Use the new attack in controller_user.go StartAttack method


*/

// Attacker interface defines the methods that an attacker should implement
type Attacker interface {
	DelayAllPacketsBy(int) error          // delay messages by delay milliseconds
	LossPercentagePackets(int) error      // loss percentage of packets
	DuplicatePercentagePackets(int) error // duplicate percentage of packets
	ReorderPercentagePackets(int) error   // reorder percentage of packets
	CorruptPercentagePackets(int) error   // corrupt percentage of packets
	Halt() error                          // halt the process
	Reset() error                         // reset the process
	Kill() error                          // kill the process
	BufferAllMessages() error             // buffer all messages
	AllowMessages(int) error              // allow num_messages messages
	CorruptDB() error                     // corrupt the internal database
}

type OperationTypes struct {
	DelayAllPacketsBy          int32
	LossPercentagePackets      int32
	DuplicatePercentagePackets int32
	ReorderPercentagePackets   int32
	CorruptPercentagePackets   int32
	Halt                       int32
	Reset                      int32
	Kill                       int32
	BufferAllMessages          int32
	AllowMessages              int32
	CorruptDB                  int32
}

func NewOperationTypes() OperationTypes {
	return OperationTypes{
		DelayAllPacketsBy:          1,
		LossPercentagePackets:      2,
		DuplicatePercentagePackets: 3,
		ReorderPercentagePackets:   4,
		CorruptPercentagePackets:   5,
		Halt:                       6,
		Reset:                      7,
		Kill:                       8,
		BufferAllMessages:          9,
		AllowMessages:              10,
		CorruptDB:                  11,
	}
}
