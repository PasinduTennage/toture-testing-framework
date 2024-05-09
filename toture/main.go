package main

import (
	"flag"
	"time"
)

func main() {
	replicaName := flag.String("rName", "replica", "unix process name of the replica")
	debugOn := flag.Bool("debugOn", false, "turn on/off debug, turn off when benchmarking")
	numReplicas := flag.Int("numReplicas", 5, "total number of replicas")
	testTime := flag.Int("testTime", 60, "total time of the test in seconds")
	benchmark := flag.String("benchmark", "none", "type of attack -- none, delay, loss, duplicate, reorder, corrupt, partition")
	attacker := flag.String("attacker", "localNetEm", "localNetEm / remoteNetEm")

	numThreshold := flag.Int("numThreshold", 2, "number of replicas attacked at the same time")
	viewChangeTime := flag.Int("viewChangeTime", 10, "time to wait for view change in milli seconds")
	epochTime := flag.Int("epochTime", 2, "maximum attack duration before changing the target in seconds")
	setup := flag.String("setup", "none", "for distributed scenarios, the IP of each node")
	//
	////
	flag.Parse()
	time.Sleep(5 * time.Second) // this wait time is for the connections to establish

}
