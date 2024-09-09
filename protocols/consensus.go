package protocols

import (
	"toture-test/consenbench/common"
	"toture-test/util"
)

type ConsensusOptions struct {
	Option map[string]string
}

type Consensus interface {
	Bootstrap(nodes []*common.Node, options ConsensusOptions) error
	Get_Performance_Stats(nodes []*common.Node, performance util.Performance, crashed_nodes []*common.Node, options ConsensusOptions) util.Performance
	ExtractOptions(string) ConsensusOptions
}
