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
	name := flag.Int64("name", 1, "name of the torture")
	debugOn := flag.Bool("debugOn", false, "true / false")
	debugLevel := flag.Int("debugLevel", 1, "debug level")
	isController := flag.Bool("isController", false, "true for controller, false for client")

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
		cl := torture.NewClient(int(*name), *cfg, *debugOn, *debugLevel)
		cl.NetworkInit()
		cl.SetAttacker(attacker_impl.NewLocalNetEmAttacker(int(*name), *debugOn, *debugLevel, nil, *cfg))
		cl.ConnectToController()
		/*to avoid exiting the main thread*/
		for true {
			time.Sleep(10 * time.Second)
		}
	}

}
