package main

import (
	"log"
	"os"

	"github.com/yeqown/arp"
)

func main() {
	local := arp.NewLocal("en0")
	resps, err := arp.SendPacket(local.Iface, local.IPNet, "192.168.1.4")

	if err != nil {
		log.Printf("could not send packet: %v", err)
		os.Exit(-1)
	}
	for _, resp := range resps {
		println(resp.String())
	}
}
