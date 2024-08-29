package consensus

import "toture-test/util"

type Consensus interface {
	Bootstrap(nodes []util.Node)
	Get_Performance_Stats(nodes []util.Node, performance util.Performance, crashed_nodes []int)
}
