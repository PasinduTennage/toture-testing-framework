package cmd

// high level attack interface
func (lna *LocalNetEmAttacker) Run() error {
	lna.Start()

	lna.End()
	return nil
}
