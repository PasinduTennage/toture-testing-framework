package dummy

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// print the memory usage of the process

func (pr *Proxy) printMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("\nTotal Memory Allocs = %v MiB\n", bToMb(m.TotalAlloc))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

// print the CPU usage of the process

func (pr *Proxy) printCPUUsage() {
	// Using 'ps' command to get CPU usage
	cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", os.Getpid()), "-o", "%cpu")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error getting CPU usage:", err)
		return
	}
	fmt.Printf("CPU Usage: %s\n", string(out))
}

// write state

func (pr *Proxy) WriteStat() {
	go func() {
		for true {
			pr.printMemoryUsage()
			pr.printCPUUsage()
			time.Sleep(10 * time.Second)
		}
	}()
}
