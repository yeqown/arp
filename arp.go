package arp

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// PacketResponse ... includes ARP packet reply: ip and MAC addr
type PacketResponse struct {
	IP  net.IP
	MAC net.HardwareAddr
}

func (r PacketResponse) String() string {
	return fmt.Sprintf("ip: %s, mac: %s", r.IP.String(), r.MAC.String())
}

// SendPacket ...
func SendPacket(iface *net.Interface, srcIPNet *net.IPNet, dstIPString string) ([]PacketResponse, error) {
	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("could not open pcap handle: %v", err)
	}
	defer handle.Close()

	arps := make(chan PacketResponse, 10)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go readARP(ctx, handle, iface.HardwareAddr, arps)
	dstIP := parseIPString(dstIPString)
	if err := writeARP(handle, iface.HardwareAddr, srcIPNet.IP, dstIP); err != nil {
		cancel()
		return nil, fmt.Errorf("error writing packets on %v: %v", iface.Name, err)
	}

	resps := []PacketResponse{}
	for arp := range arps {
		resps = append(resps, arp)
	}

	return resps, nil
}

func writeARP(handle *pcap.Handle, srcMAC net.HardwareAddr, srcIP, dstIP net.IP) error {
	println(net.IP(srcIP).String(), net.IP(dstIP).String())
	// Set up all the layers' fields we can.
	eth := layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(srcMAC),
		SourceProtAddress: []byte(srcIP),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(dstIP),
	}
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	gopacket.SerializeLayers(buf, opts, &eth, &arp)
	return handle.WritePacketData(buf.Bytes())
}

func readARP(ctx context.Context, handle *pcap.Handle, srcMAC net.HardwareAddr, arps chan PacketResponse) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	defer close(arps)
	for {
		var packet gopacket.Packet
		select {
		case <-ctx.Done():
			println("timeout done")
			return
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			log.Printf("packet ip: %s, mac: %s",
				net.IP(arp.SourceProtAddress).String(),
				net.HardwareAddr(arp.SourceHwAddress).String(),
			)
			if arp.Operation != layers.ARPReply || bytes.Equal([]byte(srcMAC), arp.SourceHwAddress) {
				// This is a packet I sent.
				continue
			}
			arpResp := PacketResponse{
				IP:  net.IP(arp.SourceProtAddress),
				MAC: net.HardwareAddr(arp.SourceHwAddress),
			}
			arps <- arpResp
			log.Println(arpResp.String())
		}
	}
}
