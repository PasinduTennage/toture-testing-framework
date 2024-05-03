package cmd

type RemoteNetEmAttacker struct {
}

func (lna *RemoteNetEmAttacker) Delay(pId int, delay int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) ResetDelay(pId int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) Loss(pId int, lossRate int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) ResetLoss(pId int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) Duplicate(pId int, duplicateRate int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) ResetDuplicate(pId int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) Reorder(pId int, reorderRate int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) ResetReorder(pId int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) Corrupt(pId int, corruptRate int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) ResetCorrupt(pId int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) Halt(pId int) error {
	return nil
}

func (lna *RemoteNetEmAttacker) ResetHalt(pId int) error {
	return nil
}
