package common

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"toture-test/util"
)

type Network struct {
	Id                  int
	ListenAddress       string
	IncomingConnections map[int]*bufio.Reader
	OutgoingConnections map[int]*bufio.Writer
	OutMutex            map[int]*sync.Mutex
	OutChan             chan *RPCPairPeer
	InputChan           chan *RPCPairPeer
	RemoteAddresses     map[int]string
	logger              *util.Logger
	rpcTable            map[uint8]*RPCPair // map each RPC type (message type) to its unique number
	messageCodes        MessageCode
}

type NetworkConfig struct {
	ListenAddress   string
	RemoteAddresses map[int]string // to connect to
}

func NewNetwork(Id int, config *NetworkConfig, outChan chan *RPCPairPeer, inChan chan *RPCPairPeer, logger *util.Logger) *Network {
	return &Network{
		Id:                  Id,
		ListenAddress:       config.ListenAddress,
		IncomingConnections: make(map[int]*bufio.Reader),
		OutgoingConnections: make(map[int]*bufio.Writer),
		OutChan:             outChan,
		OutMutex:            make(map[int]*sync.Mutex),
		InputChan:           inChan,
		RemoteAddresses:     config.RemoteAddresses,
		logger:              logger,
		rpcTable:            make(map[uint8]*RPCPair),
		messageCodes:        GetRPCCodes(),
	}
}

func (n *Network) RegisterRPC(msgObj Serializable, code uint8) {
	n.rpcTable[code] = &RPCPair{Code: code, Obj: msgObj}
}

// connect to all remote nodes and then return true

func (n *Network) ConnectRemotes() error {
	for id, address := range n.RemoteAddresses {
		var b [4]byte
		bs := b[:4]

		for true {
			conn, err := net.Dial("tcp", address)
			if err == nil {
				n.OutgoingConnections[id] = bufio.NewWriter(conn)
				binary.LittleEndian.PutUint16(bs, uint16(n.Id))
				_, err := conn.Write(bs)
				if err != nil {
					panic("Error while connecting to replica " + strconv.Itoa(int(id)))
				}
				break
			} else {
				n.logger.Debug("Error while connecting to "+strconv.Itoa(int(id))+" "+err.Error(), 3)
			}
		}
	}

	return nil
}

// listen to self.ListenAddress until all expected peers are connected

func (n *Network) Listen() error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		counter := 0
		var b [4]byte
		bs := b[:4]
		Listener, err_ := net.Listen("tcp", n.ListenAddress)
		if err_ != nil {
			panic("Error while listening to incoming connections")
		}
		for counter < len(n.RemoteAddresses) {

			conn, err := Listener.Accept()
			if err != nil {
				panic(err.Error() + fmt.Sprintf("%v", err.Error()))
			}
			if _, err := io.ReadFull(conn, bs); err != nil {
				panic(err.Error() + fmt.Sprintf("%v", err.Error()))
			}
			id := int(binary.LittleEndian.Uint16(bs))

			n.IncomingConnections[id] = bufio.NewReader(conn)
			go n.HandleReadStream(n.IncomingConnections[id], id)
			counter++
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}

// read from reader and put in to self.ListenChan

func (n *Network) HandleReadStream(reader *bufio.Reader, id int) error {
	var msgType uint8
	var err error = nil

	for true {
		if msgType, err = reader.ReadByte(); err != nil {
			n.logger.Debug("Error while reading message code: connection broken from "+strconv.Itoa(int(id))+fmt.Sprintf(" %v", err.Error()), 3)
			return err
		}
		if rpair, present := n.rpcTable[msgType]; present {
			obj := rpair.Obj.New()
			if err = obj.Unmarshal(reader); err != nil {
				n.logger.Debug("Error while unmarshalling from "+strconv.Itoa(int(id))+fmt.Sprintf(" %v", err.Error()), 3)
				return err
			}
			n.OutChan <- &RPCPairPeer{
				RpcPair: &RPCPair{
					Code: msgType,
					Obj:  obj,
				},
				Peer: id,
			}
			n.logger.Debug("Pushed a message from "+strconv.Itoa(id), 0)
		} else {
			n.logger.Debug("Error received unknown message type from "+strconv.Itoa(id), 3)
			return nil
		}
	}
	return nil
}

// send to peer

func (n *Network) Send(rpc *RPCPairPeer) error {
	w := n.OutgoingConnections[rpc.Peer]
	if w == nil {
		panic("replica not found" + strconv.Itoa(rpc.Peer))
	}
	n.OutMutex[rpc.Peer].Lock()
	err := w.WriteByte(rpc.RpcPair.Code)
	if err != nil {
		n.logger.Debug("Error writing message code byte:"+err.Error(), 0)
		n.OutMutex[rpc.Peer].Unlock()
		return nil
	}
	err = rpc.RpcPair.Obj.Marshal(w)
	if err != nil {
		n.logger.Debug("Error while marshalling:"+err.Error(), 0)
		n.OutMutex[rpc.Peer].Unlock()
		return nil
	}
	err = w.Flush()
	if err != nil {
		n.logger.Debug("Error while flushing:"+err.Error(), 0)
		n.OutMutex[rpc.Peer].Unlock()
		return nil
	}
	n.OutMutex[rpc.Peer].Unlock()
	n.logger.Debug("Internal sent message to "+strconv.Itoa(rpc.Peer), 0)
	return nil
}
