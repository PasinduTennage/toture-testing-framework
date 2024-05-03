package cmd

type LocalNetEmAttacker struct {
}

func (lna *LocalNetEmAttacker) Delay(pId int, delay int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetDelay(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Loss(pId int, lossRate int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetLoss(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Duplicate(pId int, duplicateRate int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetDuplicate(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Reorder(pId int, reorderRate int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetReorder(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Corrupt(pId int, corruptRate int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetCorrupt(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) Halt(pId int) error {
	return nil
}

func (lna *LocalNetEmAttacker) ResetHalt(pId int) error {
	return nil
}
