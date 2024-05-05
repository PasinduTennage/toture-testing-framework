package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// get the process ids with the name
func GetProcessIds(name string) ([]int, error) {
	// Execute 'pgrep' command to find process IDs by name
	cmd := exec.Command("pgrep", name)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split output into lines and parse PIDs
	lines := strings.Split(string(output), "\n")
	var pids []int
	for _, line := range lines {
		if line != "" {
			pid, convErr := strconv.Atoi(line)
			if convErr != nil {
				return nil, convErr
			}
			pids = append(pids, pid)
		}
	}

	return pids, nil
}

// run a set of dummy threads with the given port, and return the channel to stop the processes
func RunDummyThreads(ports []int) ([]chan bool, int) {

	chans := []chan bool{}

	for i := 0; i < len(ports); i++ {
		nC := make(chan bool)
		chans = append(chans, nC)
		go func(fin chan bool, p int) {
			Listener, err := net.Listen("tcp", "0.0.0.0"+":"+strconv.Itoa(p))
			if err != nil {
				panic(err.Error())
				return
			}
			fmt.Printf("Listening on port %d\n", p)
			for true {
				select {
				case <-fin:
					if Listener != nil {
						Listener.Close()
					}
					fmt.Printf("closed port %d\n", p)
					return
				default:
					if Listener != nil {

					}
					//fmt.Printf("running dummy thread on port %d\n", p)
					break
				}
			}
		}(nC, ports[i])
	}
	fmt.Printf("pid: %v\n", os.Getpid())
	return chans, os.Getpid()
}
