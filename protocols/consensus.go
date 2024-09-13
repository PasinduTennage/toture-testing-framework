package protocols

import (
	"toture-test/consenbench/common"
	"toture-test/util"
)

type ConsensusOptions struct {
	Option map[string]string
}

type Consensus interface {
	CopyConsensus(nodes []*common.Node) error
	Bootstrap(nodes []*common.Node, duration int, result chan util.Performance, bootstrap_complete chan bool)
	ExtractOptions(string)
}
