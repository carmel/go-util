package proxy

import (
	"flag"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	var iface = flag.String("i", "eth0", "Interface to get packets from")
	var filter = flag.String("f", "", "BPF filter")
	var snaplen = flag.Int("s", 16<<10, "SnapLen for pcap packet capture")
	var logAllPackets = flag.Bool("v", false, "Logs every packet in great detail")

	flag.Parse()

	handle, err := pcap.OpenLive(*iface, int32(*snaplen), true, pcap.BlockForever)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	defer handle.Close()

	if *filter != "" {
		if err := handle.SetBPFFilter(*filter); err != nil {
			panic(err)
		}
	}

	//Create a new PacketDataSource
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	//Packets returns a channel of packets
	in := src.Packets()

	for {
		// var packet gopacket.Packet
		packet := <-in
		if *logAllPackets {
			log.Println(packet)
		}
		DstIpProxy(handle, packet)
	}
}

func DstIpProxy(handle *pcap.Handle, packet gopacket.Packet) {
	// dstIp, err := net.LookupHost(domain)
	// if err != nil {
	// 	return nil, err
	// }

	// get all layers
	allLayers := packet.Layers()
	// replaceLayer := make([]gopacket.SerializableLayer, 0, len(allLayers))
	var err error
	for _, v := range allLayers {
		switch v.LayerType() {
		case layers.LayerTypeIPv4:
			ipLayer := v.(*layers.IPv4)
			ipLayer.DstIP = net.IPv4(127, 0, 0, 1)

			buf := gopacket.NewSerializeBuffer()
			err = ipLayer.SerializeTo(buf, gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true})
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
			// send the packet
			err = handle.WritePacketData(buf.Bytes())
			if err != nil {
				log.Printf("Error: %s\n", err)
			}

		case layers.LayerTypeIPv6:
			ipLayer := v.(*layers.IPv6)
			ipLayer.DstIP = net.IPv6loopback

			buf := gopacket.NewSerializeBuffer()
			err = ipLayer.SerializeTo(buf, gopacket.SerializeOptions{ComputeChecksums: true, FixLengths: true})
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
			// send the packet
			err = handle.WritePacketData(buf.Bytes())
			if err != nil {
				log.Printf("Error: %s\n", err)
			}
		}

		// if l, ok := v.(gopacket.SerializableLayer); ok {
		// 	replaceLayer = append(replaceLayer, l)
		// }
	}

	// buffer := gopacket.NewSerializeBuffer()
	// err := gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{}, replaceLayer...)
	// if err != nil {
	// 	return nil, err
	// }
	// p := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
	// if p.ErrorLayer() != nil {
	// 	return nil, err
	// }
}
