package consensus

import (
	"toture-test/consenbench/common"
	"toture-test/protocols"
	"toture-test/util"
)

type Baxos struct {
}

func NewBaxos() *Baxos {
	return &Baxos{}
}

func (ba *Baxos) CopyConsensus(nodes []*common.Node, options protocols.ConsensusOptions) error {
	return nil
}

func (ba *Baxos) Bootstrap(nodes []*common.Node, options protocols.ConsensusOptions) util.Performance {
	return util.Performance{}
}

func (ba *Baxos) ExtractOptions(string) protocols.ConsensusOptions {
	return protocols.ConsensusOptions{}
}
