package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device      string = "eth0"
	snapshotLen int32  = 1024
	promiscuous bool   = false
	err         error
	timeout     time.Duration = 30 * time.Second
	handle      *pcap.Handle
)

func GoPacketAll() {

	var cardname = "eth0"
	fmt.Printf("pcap OpenLive %s....\n", cardname)
	handle, err := pcap.OpenLive(cardname, 65536, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Do something with a packet here.
		//fmt.Println(packet)
		// 	PACKET: 106 bytes, wire length 106 cap length 106 @ 2017-06-21 13:23:16.203689 +0800 CST
		// - Layer 1 (14 bytes) = Ethernet	{Contents=[..14..] Payload=[..92..] SrcMAC=00:22:a1:03:e6:e0 DstMAC=d8:9d:67:26:d8:38 EthernetType=IPv4 Length=0}
		// - Layer 2 (20 bytes) = IPv4	{Contents=[..20..] Payload=[..72..] Version=4 IHL=5 TOS=0 Length=92 Id=28002 Flags=DF FragOffset=0 TTL=56 Protocol=UDP Checksum=37348 SrcIP=59.37.51.11 DstIP=59.37.153.245 Options=[] Padding=[]}
		// - Layer 3 (08 bytes) = UDP	{Contents=[..8..] Payload=[..64..] SrcPort=15978 DstPort=8666 Length=72 Checksum=53212}
		// - Layer 4 (64 bytes) = Payload	64 byte(s)

		printPacketInfo(packet)

	}
}

func printPacketInfo(packet gopacket.Packet) {
	// Let's see if the packet is an ethernet packet
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		fmt.Println("Ethernet layer detected.")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Println("Source MAC: ", ethernetPacket.SrcMAC)
		fmt.Println("Destination MAC: ", ethernetPacket.DstMAC)
		// Ethernet type is typically IPv4 but could be ARP or other
		fmt.Println("Ethernet type: ", ethernetPacket.EthernetType)
		fmt.Println()
	}

	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		fmt.Println("Protocol: ", ip.Protocol)
		fmt.Println()
	}

	// Let's see if the packet is TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)

		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		fmt.Printf("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
		fmt.Println("Sequence number: ", tcp.Seq)
		fmt.Println()
	}

	// Iterate over all layers, printing out each layer type
	fmt.Println("All packet layers:")
	for _, layer := range packet.Layers() {
		fmt.Println("- ", layer.LayerType())
	}

	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		fmt.Println("Application layer/Payload found.")
		//fmt.Printf("%s\n", applicationLayer.Payload())

		// Search for a string inside the payload
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
			fmt.Println("HTTP found!")
		}
	}

	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}
}

func openTargetLive(cardname string) (handle *pcap.Handle, _ error) {

	fmt.Printf("pcap OpenLive %s....\n", cardname)
	handle, err := pcap.OpenLive(cardname, 65536, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer handle.Close()
	return handle, nil
}

func setFilter(handle *pcap.Handle) {
	// Set filter
	var filter string = "tcp and port 80"
	var err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing TCP port 80 packets.")
}

func FindDevices() {
	fmt.Println("----------Find all devices---------\n ")
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// Print device information
	for _, device := range devices {
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
	/*- IP address:  45.33.110.101
	  - Subnet mask:  ffffff00
	  - IP address:  2600:3c01::f03c:91ff:fee5:45b6
	  - Subnet mask:  ffffffffffffffff0000000000000000
	  - IP address:  fe80::f03c:91ff:fee5:45b6
	  - Subnet mask:  ffffffffffffffff0000000000000000
	  - IP address:  127.0.0.1
	  - Subnet mask:  ff000000
	  - IP address:  ::1
	  - Subnet mask:  ffffffffffffffffffffffffffffffff
	*/
}

func packetLayer(handle *pcap.Handle) {
	//defer handle.Close()
	fmt.Printf("Create a new PacketDataSource.....\n")
	//Create a new PacketDataSource
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	//Packets returns a channel of packets
	in := src.Packets()
	fmt.Printf("start packets.....\n")
	for {
		var packet gopacket.Packet
		select {
		//case <-stop:
		//return
		case packet = <-in:
			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer == nil {
				//fmt.Printf("tcpLayer is nil.....\n")
				continue
			}
			tcp := tcpLayer.(*layers.TCP)

			layerType := tcp.SrcPort.LayerType().String()
			ack := tcp.ACK
			dstlayerType := tcp.DstPort.LayerType().String()

			fmt.Printf("layerType is: %s ack: %s dstlayerType: %s \n", layerType, ack, dstlayerType)

		}
	}

}
