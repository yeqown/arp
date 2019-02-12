package main

import (
	"log"

	"github.com/yeqown/arp"
	"github.com/yeqown/arp/v2"
)

func main() {
	local := arp.NewLocal("en0")
	if err := v2.ARPCall(local.Iface); err != nil {
		log.Fatal(err)
	}
}
