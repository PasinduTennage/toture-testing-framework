package configuration

import (
	"bufio"
	"os"
	"strings"
)

type ConsensusConfig struct {
	options map[string]string
}

func NewConsensusConfig(fname string) (*ConsensusConfig, error) {
	cfg := ConsensusConfig{
		options: make(map[string]string),
	}

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

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line
	for scanner.Scan() {
		// Append each line to the slice
		lines = append(lines, scanner.Text())
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	// Iterate over each line
	for _, line := range lines {
		// Split the line by the space character
		parts := strings.Split(line, "-")
		cfg.options[parts[0]] = parts[1]
	}

	return &cfg, nil
}
