package protocols

import (
	"toture-test/consenbench/common"
	"toture-test/util"
)

type ConsensusOptions struct {
	Option map[string]string
}

type Consensus interface {
	CopyConsensus(nodes []*common.Node, options ConsensusOptions) error
	Bootstrap(nodes []*common.Node, options ConsensusOptions) util.Performance
	ExtractOptions(string) ConsensusOptions
}
