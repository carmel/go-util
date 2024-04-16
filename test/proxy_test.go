package util

// import (
// 	"flag"
// 	"fmt"
// 	"goUtil/proxy"
// 	"log"
// 	"net"
// 	"testing"

// 	"github.com/google/gopacket"
// 	"github.com/google/gopacket/layers"
// 	"github.com/google/gopacket/pcap"
// )

// func TestLookupHost(t *testing.T) {
// 	domain := "www.tracup.com"
// 	ns, err := net.LookupHost(domain)
// 	if err != nil {
// 		fmt.Printf("error: %v, failed to parse %v\n", err, domain)
// 	}
// 	fmt.Printf("%s\n", ns)
// }

// func TestDstipProxy(t *testing.T) {
// 	var iface = flag.String("i", "en1", "Interface to get packets from")
// 	var filter = flag.String("f", "dst host 120.76.153.29", "BPF filter") // dst host 192.9.200.1 or dst host 192.9.200.2
// 	var snaplen = flag.Int("s", 16<<10, "SnapLen for pcap packet capture")
// 	var logAllPackets = flag.Bool("v", true, "Logs every packet in great detail")

// 	flag.Parse()
// 	handle, err := pcap.OpenLive(*iface, int32(*snaplen), true, pcap.BlockForever)
// 	if err != nil {
// 		log.Printf("Error: %s\n", err)
// 		return
// 	}
// 	defer handle.Close()

// 	if *filter != "" {
// 		if err := handle.SetBPFFilter(*filter); err != nil {
// 			panic(err)
// 		}
// 	}

// 	//Create a new PacketDataSource
// 	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
// 	//Packets returns a channel of packets
// 	in := src.Packets()

// 	for {
// 		// var packet gopacket.Packet
// 		packet := <-in
// 		if *logAllPackets {
// 			log.Println(packet)
// 		}
// 		proxy.DstIpProxy(handle, packet)
// 	}
// }
