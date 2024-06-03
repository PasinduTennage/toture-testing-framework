package main

import (
	"flag"
	"time"
	"toture-test/torture/attacker_impl"
	"toture-test/torture/configuration"
	"toture-test/torture/controller_frontend"
	torture "toture-test/torture/torture/src"
)

// this file defines the main routine of Torture (client or controller)

func main() {
	configFile := flag.String("config", "torture/configuration/local-config.cfg", "configuration file")
	replicaConfigFile := flag.String("replicaConfig", "torture/configuration/consensus_config/11.cfg", "consensus configuration file (contains consensus replica information)")
	name := flag.Int64("name", 11, "name of the torture")
	debugOn := flag.Bool("debugOn", false, "true / false")
	debugLevel := flag.Int("debugLevel", 1, "debug level")
	isController := flag.Bool("isController", false, "true for controller, false for client")
	attacker := flag.String("attacker", "localNetEm", "localNetEm / remoteNetEm")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		panic(err.Error())
	}

	if *isController {
		c := torture.NewController(int(*name), *cfg, *debugOn, *debugLevel)
		c.Run()
		c.NetworkInit()
		c.ConnectToClients()
		time.Sleep(10 * time.Second)
		nodes := torture.CreateNodes(*cfg, c)
		controller_frontend.StartAttack(nodes)
	} else {
		consensus_config, err := configuration.NewConsensusConfig(*replicaConfigFile)
		if err != nil {
			panic(err.Error())
		}
		cl := torture.NewClient(int(*name), *cfg, *debugOn, *debugLevel)
		cl.NetworkInit()
		if *attacker == "localNetEm" {
			cl.SetAttacker(attacker_impl.NewLocalNetEmAttacker(int(*name), *debugOn, *debugLevel, *cfg, *consensus_config, cl))
		} else if *attacker == "remoteNetEm" {
			cl.SetAttacker(attacker_impl.NewRemoteNetEmAttacker(int(*name), *debugOn, *debugLevel, *cfg, *consensus_config, cl))
		} else {
			panic("invalid")
		}
		cl.ConnectToController()
		/*to avoid exiting the main thread*/
		for true {
			time.Sleep(10 * time.Second)
		}
	}

}
