package consensus

import (
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type Baxos struct {
}

func (ba *Baxos) Bootstrap(nodes []*common.Node, options protocols.ConsensusOptions) error {
	return nil
	// logic to bootstrap
}

func (ba *Baxos) Get_Performance_Stats(nodes []*common.Node, performance util.Performance, crashed_nodes []*common.Node, options protocols.ConsensusOptions) util.Performance {
	return performance // get performance stats from the nodes except the crashed and update the performance object
}

func (ba *Baxos) ExtractOptions(optionsFile string) protocols.ConsensusOptions {
	return protocols.ConsensusOptions{}
}
