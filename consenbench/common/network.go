package common

import (
	"bufio"
	"sync"
)

type Network struct {
	ListenAddress       string
	IncomingConnections map[int]*bufio.Reader
	OutgoingConnections map[int]*bufio.Writer
	OutChan             chan OutgoingRPC
	OutMutex            map[int]*sync.Mutex
	InputChan           chan OutgoingRPC
	RemoteAddresses     []string
}

type NetworkConfig struct {
	ListenAddress   string
	RemoteAddresses []string // to connect to
}

func NewNetwork(config *NetworkConfig, outChan chan OutgoingRPC, inChan chan OutgoingRPC) *Network {
	return &Network{
		ListenAddress:       config.ListenAddress,
		IncomingConnections: make(map[int]*bufio.Reader),
		OutgoingConnections: make(map[int]*bufio.Writer),
		OutChan:             outChan,
		OutMutex:            make(map[int]*sync.Mutex),
		InputChan:           inChan,
		RemoteAddresses:     config.RemoteAddresses,
	}
}

func (n *Network) ConnectRemote() error {
	// connect to all remote nodes
	return nil
}

func (n *Network) Listen() error {
	// listen to self.ListenAddress
	return nil
}

func (n *Network) HandleReadStream(reader *bufio.Reader) error {
	// read from reader and put in to self.ListenChan
	return nil
}

func (n *Network) Send(peer int) error {
	// send to peer
	return nil
}
