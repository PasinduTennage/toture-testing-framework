package cmd

import (
	"bufio"
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

// run a set of dummy threads with the given port, and return the channel to stop the threads

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
						//fmt.Printf("running dummy thread on port %d\n", p)
					}
					break
				}
			}
		}(nC, ports[i])
	}
	fmt.Printf("pid: %v\n", os.Getpid())
	return chans, os.Getpid()
}

// given a port return the process id

func GetProcessID(port int) int {
	// Run the lsof command to get the process ID using the specified port
	out, err := exec.Command("lsof", "-ti", ":"+strconv.Itoa(port), "-s", "TCP:LISTEN").Output()
	if err != nil {
		panic(err.Error())
	}

	// Split the output into lines
	lines := strings.Split(string(out), "\n")

	if len(lines) >= 1 {
		fields := strings.Fields(lines[0])
		if len(fields) >= 1 {
			pid, _ := strconv.Atoi(fields[0])
			return pid
		}
	}

	// If no process ID found, return -1
	return -1
}

// NewConfig loads [][]ports from a file

func NewConfig(fname string) ([][]int, error) {
	ports := make([][]int, 0)
	file, err := os.Open(fname)
	if err != nil {
		panic(err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err.Error())
		}
	}(file)

	var lines []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	for _, line := range lines {
		parts := strings.Split(line, " ")
		new_ports := convertToIntArr(parts[2:])
		ports = append(ports, new_ports)
	}

	return ports, nil
}

// convert string array to int array
func convertToIntArr(i []string) []int {
	var res []int
	for _, v := range i {
		val, err := strconv.Atoi(v)
		if err != nil {
			panic(err.Error())
		}
		res = append(res, val)
	}
	return res
}
