package common

import "bufio"

type Network struct {
	ListenAddress       string
	IncomingConnections map[int]*bufio.Reader
	OutgoingConnections map[int]*bufio.Writer
	ListenChan          chan interface{}
}

type NetworkConfig struct {
	ListenAddress   string
	RemoteAddresses []string // to connect to
}

func NewNetwork(config *NetworkConfig, listenChan chan interface{}) *Network {
	return &Network{
		ListenAddress:       config.ListenAddress,
		IncomingConnections: make(map[int]*bufio.Reader),
		OutgoingConnections: make(map[int]*bufio.Writer),
		ListenChan:          listenChan,
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
