package controller

type AttackNode struct {
	Id         int
	Controller *Controller
}

func (n *AttackNode) Kill() {

}

func (n *AttackNode) Slowdown() {

}

func (n *AttackNode) Pause() {

}

func (n *AttackNode) Continue() {

}

func (n *AttackNode) IsLeader() {

}
func (n *AttackNode) SetSkew() {

}
func (n *AttackNode) SetDrift() {

}

type AttackLink struct {
	Id_sender   int
	Id_reciever int
	Controller  *Controller
}

func (a *AttackLink) SetStatus() {

}

func (a *AttackLink) SetDelay() {

}

func (a *AttackLink) SetLoss() {

}

func (a *AttackLink) SetBandwidth() {

}

type Attacker interface {
	Attack([]AttackNode, [][]AttackLink, int) // duration
}
