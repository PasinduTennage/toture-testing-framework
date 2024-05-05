package main

import (
	"flag"
	"net"
	"strings"
	"time"
)

func main() {
	ports := flag.String("ports", "10000,100001", "set of open ports for this process")

	flag.Parse()

	// convert the ports string to a slice of integers
	portSlice := strings.Split(*ports, ",")

	for i := 0; i < len(portSlice); i++ {
		go func(port string) {
			net.Listen("tcp", "0.0.0.0:"+port)
		}(portSlice[i])
	}
	for true {
		time.Sleep(100 * time.Second)
	}

}
