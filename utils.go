package arp

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"
)

func parseIPString(s string) net.IP {
	var byts []byte
	for _, i := range strings.Split(s, ".") {
		v, _ := strconv.Atoi(i)
		byts = append(byts, uint8(v))
	}
	return byts
}

// ips is a simple and not very good method for getting all IPv4 addresses from a
// net.IPNet.  It returns all IPs it can over the channel it sends back, closing
// the channel when done.
func ips(addr *net.IPNet) (out []net.IP) {
	num := binary.BigEndian.Uint32([]byte(addr.IP))
	mask := binary.BigEndian.Uint32([]byte(addr.Mask))
	num &= mask
	for mask < 0xffffffff {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], num)
		out = append(out, net.IP(buf[:]))
		mask++
		num++
	}
	return
}
