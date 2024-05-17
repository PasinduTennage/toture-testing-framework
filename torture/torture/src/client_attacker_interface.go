package torture

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
}
