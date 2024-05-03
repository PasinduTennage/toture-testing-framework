package main

import (
	"flag"
	"time"
)

func main() {
	replicaName := flag.String("rName", "replica", "unix process name of the replica")
	debugOn := flag.Bool("debugOn", false, "turn on/off debug, turn off when benchmarking")
	numReplicas := flag.Int("numReplicas", 5, "total number of replicas")
	numThreshold := flag.Int("numThreshold", 2, "number of replicas attacked at the same time")
	testTime := flag.Int("testTime", 60, "total time of the test in seconds")
	viewChangeTime := flag.Int("viewChangeTime", 10, "time to wait for view change in milli seconds")
	benchmark := flag.String("benchmark", "none", "type of attack -- none, delay, loss, duplicate, reorder, corrupt, partition")
	epochTime := flag.Int("epochTime", 2, "maximum attack duration before changing the target in seconds")
	configuration := flag.String("configuration", "none", "for distributed scenarios, the IP of each node")
	mode := flag.String("mode", "local", "local / remote")

	//
	flag.Parse()
	time.Sleep(5 * time.Second) // this wait time is for the connections to establish

}
