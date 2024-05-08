package cmd

func (lna *LocalNetEmAttacker) Run() error {
	lna.Start()

	lna.End()
	return nil
}
