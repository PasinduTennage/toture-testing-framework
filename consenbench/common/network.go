package common

import (
	"bufio"
	"sync"
)

type Network struct {
	ListenAddress       string
	IncomingConnections map[int]*bufio.Reader
	OutgoingConnections map[int]*bufio.Writer
	OutChan             chan *RPCPairPeer
	OutMutex            map[int]*sync.Mutex
	InputChan           chan *RPCPairPeer
	RemoteAddresses     []string
}

type NetworkConfig struct {
	ListenAddress   string
	RemoteAddresses []string // to connect to
}

func NewNetwork(config *NetworkConfig, outChan chan *RPCPairPeer, inChan chan *RPCPairPeer) *Network {
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
	// connect to all remote nodes and then return true
	return nil
}

func (n *Network) Listen() error {
	// listen to self.ListenAddress until all expected peers are connected
	return nil
}

func (n *Network) HandleReadStream(reader *bufio.Reader, peer int) error {
	// read from reader and put in to self.ListenChan
	return nil
}

func (n *Network) Send(rpc *RPCPairPeer) error {
	// send to peer
	return nil
}
