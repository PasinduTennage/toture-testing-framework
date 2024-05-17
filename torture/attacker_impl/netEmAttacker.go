package attacker_impl

type LocalNetEmAttacker struct {
	debugOn    bool
	debugLevel int
}

// NewLocalNetEmAttacker creates a new LocalNetEmAttacker

func NewLocalNetEmAttacker() *LocalNetEmAttacker {
	return nil
}

func (l *LocalNetEmAttacker) DelayAllPacketsBy(int) error {
	return nil
}

func (l *LocalNetEmAttacker) LossPercentagePackets(int) error {
	return nil
}

func (l *LocalNetEmAttacker) DuplicatePercentagePackets(int) error {
	return nil
}

func (l *LocalNetEmAttacker) ReorderPercentagePackets(int) error {
	return nil
}

func (l *LocalNetEmAttacker) CorruptPercentagePackets(int) error {
	return nil
}

func (l *LocalNetEmAttacker) Halt() error {
	return nil
}

func (l *LocalNetEmAttacker) Reset() error {
	return nil
}

func (l *LocalNetEmAttacker) Kill() error {
	return nil
}

func (l *LocalNetEmAttacker) BufferAllMessages() error {
	return nil
}

func (l *LocalNetEmAttacker) AllowMessages(int) error {
	return nil
}
