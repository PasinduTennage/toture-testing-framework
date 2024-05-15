package main

import (
	"flag"
	"time"
	"toture-test/torture/configuration"
)

// this file defines the main routine of Dummy, which takes input arguments from the command line

func main() {
	configFile := flag.String("config", "dummy/configuration/local-config.cfg", "configuration file")
	name := flag.Int64("name", 1, "name of the torture")
	debugOn := flag.Bool("debugOn", false, "true / false")
	debugLevel := flag.Int("debugLevel", 1, "debug level")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		panic(err.Error())
	}

	proxyInstance := dummy.NewProxy(*name, *cfg, *debugOn, *debugLevel)

	proxyInstance.NetworkInit()
	proxyInstance.Run()
	time.Sleep(10 * time.Second)
	proxyInstance.ConnectToReplicas()
	time.Sleep(10 * time.Second)
	proxyInstance.WriteStat()
	proxyInstance.StartApplication()

	/*to avoid exiting the main thread*/
	for true {
		time.Sleep(10 * time.Second)
	}
}
