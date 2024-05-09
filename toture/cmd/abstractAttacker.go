package cmd

// Attacker interface defines the methods that an attacker should implement
type Attacker interface {
	Delay(int, int) error // delay messages by delay milliseconds

	Loss(int, int) error // drop lossRate percentage of messages

	Duplicate(int, int) error // duplicate duplicateRate percentage of messages

	Reorder(int, int) error // reorder reorderRate percentage of messages

	Corrupt(int, int) error // corrupt corruptRate percentage of messages

	Halt(int) error // halt the process

	Reset(int) error // reset the current attack

	Kill(int) error // kill the process

	Start() error // start the attack
	End() error   // start the attack

	GetPiDPortMap() map[int][]int // get the map of process id to port
}
