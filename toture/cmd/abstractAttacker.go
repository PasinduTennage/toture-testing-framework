package cmd

// Attacker interface defines the methods that an attacker should implement
type Attacker interface {
	Delay(pId int, delay int) error // delay messages by delay milliseconds

	Loss(pId int, lossRate int) error // drop lossRate percentage of messages

	Duplicate(pId int, duplicateRate int) error // duplicate duplicateRate percentage of messages

	Reorder(pId int, reorderRate int) error // reorder reorderRate percentage of messages

	Corrupt(pId int, corruptRate int) error // corrupt corruptRate percentage of messages

	Halt(pId int) error // halt the process

	Reset(pId int) error // reset the current attack

	Kill(pId int) error // kill the process

	Run() error   // start the attack
	Start() error // start the attack
	End() error   // start the attack
}
