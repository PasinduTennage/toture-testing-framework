package main

import (
	"flag"
	"fmt"
	"toture-test/toture/cmd"
)

func main() {
	replicaName := flag.String("rName", "replica", "unix process name of the replica")
	debugOn := flag.Bool("debugOn", false, "turn on/off debug, turn off when benchmarking")
	debugLevel := flag.Int("debugLevel", 0, "current debug level")
	numReplicas := flag.Int("numReplicas", 5, "total number of replicas")
	testTime := flag.Int("testTime", 60, "total time of the test in seconds")
	attacker := flag.String("attacker", "localNetEm", "localNetEm / remoteNetEm")

	numThreshold := flag.Int("numThreshold", 2, "number of replicas attacked at the same time")
	viewChangeTime := flag.Int("viewChangeTime", 10, "time to wait for view change in milli seconds")
	epochTime := flag.Int("epochTime", 2, "maximum attack duration before changing the target in seconds")

	delay := flag.Int("delay", 200, "delay in milliseconds")
	lossRate := flag.Int("lossRate", 2, "loss rate")
	duplicateRate := flag.Int("duplicateRate", 2, "duplicate rate")
	reorderRate := flag.Int("reorderRate", 2, "reorder rate")
	corruptRate := flag.Int("corruptRate", 2, "corrupt rate")

	config := flag.String("config", "toture/configuration/dummy_config.cfg", "config file for the deployment")

	flag.Parse()

	options := make(map[string]any)
	options["replicaName"] = *replicaName
	options["debugOn"] = *debugOn
	options["debugLevel"] = *debugLevel
	options["numReplicas"] = *numReplicas
	options["testTime"] = *testTime
	options["attacker"] = *attacker
	options["numThreshold"] = *numThreshold
	options["viewChangeTime"] = *viewChangeTime
	options["epochTime"] = *epochTime
	options["delay"] = *delay
	options["lossRate"] = *lossRate
	options["duplicateRate"] = *duplicateRate
	options["reorderRate"] = *reorderRate
	options["corruptRate"] = *corruptRate

	var attackerInstance cmd.Attacker

	switch *attacker {
	case "localNetEm":
		if *debugOn {
			fmt.Printf("starting localNetEm attacker\n")
		}
		ports, _ := cmd.NewConfig(*config)
		attackerInstance = cmd.NewLocalNetEmAttacker(*replicaName, ports, options, *debugOn, *debugLevel)
		break
	default:
		break
	}

	t := cmd.NewTorture(attackerInstance, options, *debugOn, *debugLevel)
	t.Run()
}
