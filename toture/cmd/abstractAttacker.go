package cmd

// Attacker interface defines the methods that an attacker should implement
type Attacker interface {
	Delay(pId int, delay int) error // delay messages by delay milliseconds
	ResetDelay(pId int) error       // reset the delay

	Loss(pId int, lossRate int) error // drop lossRate percentage of messages
	ResetLoss(pId int) error          // reset the loss

	Duplicate(pId int, duplicateRate int) error // duplicate duplicateRate percentage of messages
	ResetDuplicate(pId int) error               // reset the duplicate

	Reorder(pId int, reorderRate int) error // reorder reorderRate percentage of messages
	ResetReorder(pId int) error             // reset the reorder

	Corrupt(pId int, corruptRate int) error // corrupt corruptRate percentage of messages
	ResetCorrupt(pId int) error             // reset the corrupt

	Halt(pId int) error      // halt the process
	ResetHalt(pId int) error // reset the halt

	StartAttack() error // start the attack
	StopAttack() error  // stop the attack
}
