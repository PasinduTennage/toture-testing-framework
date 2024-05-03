package cmd

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// get the process ids with the name
func getProcessIds(name string) ([]int, error) {
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

// get the port associated with the given process IDs (we assume one port per process)
func getPorts(pIds []int) []int {
	var ports []int

	for _, pid := range pIds {
		// Run netstat command to get network connections
		cmd := exec.Command("netstat", "-anp")
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("Error running netstat command:", err)
			continue
		}

		// Parse the output to get ports associated with the given PID
		lines := strings.Split(string(out), "\n")
		line := lines[0]
		fields := strings.Fields(line)
		if len(fields) > 6 && strings.HasPrefix(fields[6], strconv.Itoa(pid)+"/") {
			s := strings.Split(fields[3], ":")
			if len(s) < 2 {
				continue
			}
			// Extract port number
			port := strings.Split(fields[3], ":")[1]
			portNum, _ := strconv.Atoi(port)
			ports = append(ports, portNum)
		}

	}

	return ports
}

// run a dummy process with the given name and open port, and return the process ID, and a channel to stop the process
func runDummyProcess(name string, port int) int {

}
