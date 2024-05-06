package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"toture-test/dummy/configuration"
	dummy "toture-test/dummy/replica/src"
)

// this file defines the main routine of QuePaxa, which takes input arguments from the command line

func main() {
	configFile := flag.String("config", "dummy/configuration/local-config.txt", "configuration file")
	name := flag.Int64("name", 1, "name of the replica")
	debugOn := flag.Bool("debugOn", false, "true / false")
	debugLevel := flag.Int("debugLevel", 1010, "debug level")

	flag.Parse()

	cfg, err := configuration.NewInstanceConfig(*configFile, *name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		os.Exit(1)
	}

	proxyInstance := dummy.NewProxy(*name, *cfg, *debugOn, *debugLevel)

	proxyInstance.NetworkInit()
	proxyInstance.Run()
	time.Sleep(10 * time.Second)
	proxyInstance.ConnectToReplicas()
	time.Sleep(10 * time.Second)
	proxyInstance.StartApplication()

	/*to avoid exiting the main thread*/
	for true {
		time.Sleep(10 * time.Second)
	}
}
