package main

import (
	"fmt"
	"log"
	"time"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var portsToCapture = map[uint16]bool{
	10000: true, 10001: true, 10002: true, 10003: true,
	20000: true, 20001: true, 20002: true, 20003: true, 20004: true,
	30000: true, 30001: true, 30002: true, 30003: true,
	40000: true, 40001: true, 40002: true, 40003: true,
}

func main() {
	nfq, err := netfilter.NewNFQueue(1, 1000000, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		log.Fatalf("could not open NFQUEUE: %v", err)
	}
	defer nfq.Close()

	packets := nfq.GetPackets()

	for packet := range packets {
		handlePacket(packet)
	}
}

func handlePacket(packet netfilter.NFPacket) {
	data := packet.Packet.Data()
	packetGopacket := gopacket.NewPacket(data, layers.LayerTypeIPv4, gopacket.Default)

	ipLayer := packetGopacket.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		packet.SetVerdict(netfilter.NF_ACCEPT)
		return
	}

	tcpLayer := packetGopacket.Layer(layers.LayerTypeTCP)
	if tcpLayer == nil {
		packet.SetVerdict(netfilter.NF_ACCEPT)
		return
	}
	tcp, _ := tcpLayer.(*layers.TCP)

	if _, ok := portsToCapture[uint16(tcp.DstPort)]; ok {
		fmt.Printf("Delaying packet: %v\n", packetGopacket)
		time.Sleep(2 * time.Millisecond)
		packet.SetVerdict(netfilter.NF_ACCEPT)
	} else {
		packet.SetVerdict(netfilter.NF_ACCEPT)
	}
}
