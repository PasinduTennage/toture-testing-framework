package attacks

import "toture-test/consenbench/controller"

type BasicAttack struct {
}

func (a *BasicAttack) Attack(nodes []controller.AttackNode, links [][]controller.AttackLink, duration int) {
}
