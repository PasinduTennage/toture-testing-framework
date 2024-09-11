package consenbench

import (
	"flag"
	"strings"
	"toture-test/consenbench/client"
	"toture-test/consenbench/controller"
)

func main() {
	// generic arguments
	is_controller := flag.Bool("is_controller", false, "is the node a controller")
	node_info_file := flag.String("node_config", "consenbench/assets/ip.yaml", "node ip configuration file")
	id := flag.Int("id", 1, "id of the node, according to node ip configuration file")
	debug_on := flag.Bool("debug_on", false, "turn on debug mode")
	debug_level := flag.Int("debug_level", 0, "debug level")

	// controller specific arguments
	attack_duration := flag.Int("attack_duration", 60, "duration of the attack in seconds")
	attacks := flag.String("attacks", "all", "set of attacks to run, seperated by commas without spaces")
	controller_operation_type := flag.String("controller_operation_type", "", "operation type of the controller: BootstrapClients/CopyConsensus/Run")
	consensus_algorithm := flag.String("consensus_algorithm", "raft", "consensus algorithm to run")

	flag.Parse()

	if *is_controller {
		options := controller.ControllerOptions{
			AttackDuration: *attack_duration,
			Attacks:        strings.Split(*attacks, ","),
			NodeInfoFile:   *node_info_file,
			DebugOn:        *debug_on,
			DebugLevel:     *debug_level,
		}
		controller := controller.NewController(*id, options)
		if *controller_operation_type == "BootstrapClients" {
			controller.BootstrapClients()
		} else if *controller_operation_type == "CopyConsensus" {
			controller.CopyConsensus(*consensus_algorithm)
		} else if *controller_operation_type == "Run" {
			controller.Run(*consensus_algorithm)
		} else {
			panic("invalid controller operation type")
		}
	} else {
		options := client.ClientOptions{
			NodeInfoFile: *node_info_file,
			DebugOn:      *debug_on,
			DebugLevel:   *debug_level,
		}
		client := client.NewClient(*id, options)
		client.Run()
	}

}
