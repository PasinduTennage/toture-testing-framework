package util

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// given a port return the process id

func GetProcessID(port string) int {
	// Run the lsof command to get the process ID using the specified port
	out, err := exec.Command("lsof", "-ti", ":"+port, "-s", "TCP:LISTEN").Output()
	for err != nil {
		time.Sleep(1 * time.Second)
		out, err = exec.Command("lsof", "-ti", ":"+port, "-s", "TCP:LISTEN").Output()
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

// runCommand runs the given command with the provided arguments -- doesn't return the command output

func RunCommand(name string, arg []string) error {
	cmd := exec.Command(name, arg...)
	if cmd.Err != nil {
		fmt.Println("Error running command " + name + " with arguments " + strings.Join(arg, " ") + " " + cmd.Err.Error() + "\n")
		return cmd.Err
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running command " + name + " with arguments " + strings.Join(arg, " ") + " " + err.Error() + "\n")
		return err
	}
	return nil
}
