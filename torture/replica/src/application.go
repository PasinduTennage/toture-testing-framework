package torture

import (
	"fmt"
	"strconv"
	"time"
)

func (pr *Proxy) handleMessage(message *Message, sender int32) {
	pr.debug("Received message: "+message.Message+" from "+string(sender), 0)
	v, ok := pr.sent.Load(message.Index)
	if ok {
		// a response to my previous message
		pr.receivedLatency = append(pr.receivedLatency, time.Now().Sub(v.(Request).sentTime).Microseconds())
		pr.sent.Delete(message.Index)

		if time.Now().Sub(pr.lastStarTime).Seconds() > 5 {
			// print the average latency
			avg := int64(0)
			for _, v := range pr.receivedLatency {
				avg += v
			}
			avg = avg / int64(len(pr.receivedLatency))
			fmt.Printf("Average latency: %d micro seconds\n", avg)
			fmt.Printf("Average throughput: %v requests per second\n\n", float64(len(pr.receivedLatency))/(time.Now().Sub(pr.lastStarTime).Seconds()))
			pr.lastStarTime = time.Now()
			pr.receivedLatency = make([]int64, 0)
		}

	} else {
		// echo the message back to the sender
		pr.sendMessage(int64(sender), message)

	}

}

// when started the application will send random messages to random nodes

func (pr *Proxy) StartApplication() {
	pr.lastStarTime = time.Now()
	node := -1
	for true {
		// select a random node that is not self
		for k, _ := range pr.addrList {
			node = int(k)
			break
		}
		rm := pr.getRandomMessage()
		pr.sendMessage(int64(node), rm)
		pr.sent.Store(rm.Index, Request{sentTime: time.Now()})
		time.Sleep(time.Duration(1) * time.Millisecond)
	}

}

// generate a random message

func (pr *Proxy) getRandomMessage() *Message {
	// generate a random Message
	pr.counter++
	return &Message{
		Index:   strconv.FormatInt(pr.name, 10) + ":" + strconv.FormatInt(pr.counter, 10),
		Message: "random message",
	}

}
