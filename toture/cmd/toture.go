package cmd

type Toture struct {
	attacker Attacker
}

func New(attacker Attacker) *Toture {
	to := Toture{
		attacker: attacker,
	}
	return &to
}

func (t *Toture) Run() error {
	return nil
}
