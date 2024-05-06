package dummy

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// start listening to the proxy tcp connection, and setup all outgoing wires

func (pr *Proxy) NetworkInit() {
	pr.waitForConnections()
}

/*
	Listen on the server ports for new connections
*/

func (pr *Proxy) waitForConnections() {

	for i := 0; i < len(pr.serverAddress); i++ {

		go func(la string) {

			var b [4]byte
			bs := b[:4]

			listener, _ := net.Listen("tcp", la)
			pr.debug("Listening to messages on "+la, 0)

			for true {
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("TCP accept error:", err)
					panic(err)
				}
				if _, err := io.ReadFull(conn, bs); err != nil {
					fmt.Println("Connection id reading error:", err)
					panic(err)
				}
				id := int32(binary.LittleEndian.Uint16(bs))
				pr.debug("Received incoming tcp connection from "+strconv.Itoa(int(id)), -1)

				go pr.connectionListener(bufio.NewReader(conn), id)
				pr.debug("Started listening to "+strconv.Itoa(int(id)), -1)
			}
		}(pr.serverAddress[i])
	}
}

/*
	listen to a given connection. Upon receiving any message, put it into the central buffer
*/

func (pr *Proxy) connectionListener(reader *bufio.Reader, id int32) {

	var err error = nil
	for true {
		obj := (&Message{}).New()
		if err = obj.Unmarshal(reader); err != nil {
			//pr.debug("Error while unmarshalling", 0)
			return
		}
		pr.incomingChan <- &ReceivedMessage{message: obj, sender: id}
	}
}

/*
	make a TCP connection to the client id
*/

func (pr *Proxy) ConnectToReplicas() {

	for id, addresses := range pr.addrList {
		for i := 0; i < len(addresses); i++ {
			var b [4]byte
			bs := b[:4]
			for true {
				conn, err := net.Dial("tcp", addresses[i])
				if err == nil {
					pr.outgoingWriters[id] = append(pr.outgoingWriters[id], bufio.NewWriter(conn))
					binary.LittleEndian.PutUint16(bs, uint16(pr.name))
					_, err := conn.Write(bs)
					if err != nil {
						//pr.debug("Error connecting to client "+strconv.Itoa(int(id)), 0)
						panic(err)
					}
					pr.debug("Started outgoing tcp connection to "+addresses[i], 0)
					break
				} else {
					time.Sleep(time.Duration(10) * time.Millisecond)
				}
			}
		}
	}

}

/*
	write a message to the wire
*/

func (pr *Proxy) sendMessage(peer int64, msg *Message) {

	pr.debug("sending message to  "+strconv.Itoa(int(peer)), 0)

	randomWriter := rand.Intn(len(pr.outgoingWriters[peer]))

	var w *bufio.Writer

	w = pr.outgoingWriters[peer][randomWriter]

	err := msg.Marshal(w)
	if err != nil {
		pr.debug("Error while marshalling", 0)
		return
	}
	err = w.Flush()
	if err != nil {
		pr.debug("Error while flushing", 0)
		return
	}
	pr.debug("sent message to  "+strconv.Itoa(int(peer)), 1)

}
