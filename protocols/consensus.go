package protocols

import "toture-test/util"

type ConsensusOptions struct {
	Option map[string]string
}

type Consensus interface {
	Bootstrap(nodes []*util.Node, options ConsensusOptions) error
	Get_Performance_Stats(nodes []*util.Node, performance util.Performance, crashed_nodes []*util.Node, options ConsensusOptions) util.Performance
	ExtractOptions() ConsensusOptions
}
