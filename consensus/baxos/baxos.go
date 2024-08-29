package consensus

import "toture-test/util"

type Baxos struct {
}

func (ba *Baxos) Bootstrap(nodes []util.Node) {
	// logic to bootstrap
}

func (ba *Baxos) Get_Performance_Stats(nodes []util.Node, performance util.Performance, crashed_nodes []int) {
	// get performance stats from the nodes except the crashed and update the performance object
}
