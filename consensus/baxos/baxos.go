package consensus

import (
	"toture-test/consensus"
	"toture-test/util"
)

type Baxos struct {
}

func (ba *Baxos) Bootstrap(nodes []util.Node, options consensus.ConsensusOptions) error {
	return nil
	// logic to bootstrap
}

func (ba *Baxos) Get_Performance_Stats(nodes []util.Node, performance util.Performance, crashed_nodes []util.Node, options consensus.ConsensusOptions) util.Performance {
	return performance // get performance stats from the nodes except the crashed and update the performance object
}

func (ba *Baxos) ExtractOptions(filepath string) consensus.ConsensusOptions {
	return consensus.ConsensusOptions{}
}
