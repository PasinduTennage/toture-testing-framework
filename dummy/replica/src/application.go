package dummy

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func (pr *Proxy) handleMessage(message *Message, sender int32) {
	pr.debug("Received message: "+message.Message, 0)
	v, ok := pr.sent.Load(message.Index)
	if ok {
		// a response to my previous message
		pr.receivedLatency = append(pr.receivedLatency, time.Now().Sub(v.(Request).sentTime).Microseconds())
		pr.sent.Delete(message.Index)

		if time.Now().Sub(pr.startTime).Seconds() > 1 {
			// print the average latency
			avg := int64(0)
			for _, v := range pr.receivedLatency {
				avg += v
			}
			avg = avg / int64(len(pr.receivedLatency))
			fmt.Printf("Average latency: %d micro seconds\n", avg)
			fmt.Printf("Average throughput: %d requests per second\n", float64(len(pr.receivedLatency))/(time.Now().Sub(pr.startTime).Seconds()))
			pr.startTime = time.Now()
			pr.receivedLatency = make([]int64, 0)
		}

	} else {
		// echo the message back to the sender
		pr.sendNewMessage(message, sender)

	}

}

func (pr *Proxy) StartApplication() {
	pr.startTime = time.Now()
	for true {
		// select a random node that is not self
		node := -1
		for k, _ := range pr.addrList {
			if k != pr.name && rand.Intn(4) == 1 {
				node = int(k)
				break
			}
		}
		rm := pr.getRandomMessage()
		pr.sendNewMessage(rm, int32(node))
		pr.sent.Store(rm.Index, Request{sentTime: time.Now()})
		time.Sleep(time.Duration(10) * time.Millisecond)
	}

}

// send a given message to given id
func (pr *Proxy) sendNewMessage(message *Message, id int32) {
	pr.debug("Sending message to "+strconv.FormatInt(int64(id), 10), 0)
	pr.sendMessage(int64(id), message)
}

// generate a random message
func (pr *Proxy) getRandomMessage() *Message {
	// generate a random message of length
	pr.counter++
	return &Message{
		Index:   strconv.FormatInt(pr.name, 10) + ":" + strconv.FormatInt(pr.counter, 10),
		Message: "random message",
	}

}
