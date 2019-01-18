// Package arp ...
package arp

import (
	"fmt"
	"log"
	"net"
	"os"
)

// NewLocal ... if netIfaceName is empty,
// arp.Local will auto find one available interface
func NewLocal(netIfaceName string) *Local {
	local := &Local{
		NetIfaceName: netIfaceName,
	}
	local.init()
	return local
}

// Local ...
type Local struct {
	NetIfaceName string
	Hostname     string           // hostname
	MAC          net.HardwareAddr // hardware addr
	IPNet        *net.IPNet       // ipnet info
	Iface        *net.Interface   // interface info
}

func (l *Local) init() {
	l.Hostname, _ = os.Hostname()
	netIfaces := make([]net.Interface, 0)

	if l.NetIfaceName == "" {
		var err error
		netIfaces, err = net.Interfaces()
		if err != nil {
			panic("cannot read local network")
		}
	} else {
		iface, err := net.InterfaceByName(l.NetIfaceName)
		if err != nil {
			panic(fmt.Sprintf("cannot read local network: %s", l.NetIfaceName))
		}
		netIfaces = append(netIfaces, *iface)
	}

	for _, netIface := range netIfaces {
		addrs, err := netIface.Addrs()
		if err != nil {
			log.Printf("[Error]: cannot load addrs from interface: %s", netIface.Name)
			continue
		}
		for _, addr := range addrs {
			// ip is not 127.x.x.x
			ip, ok := addr.(*net.IPNet)
			// println(ok, ip.IP.IsLoopback(), string(ip.IP.To4().String()))
			if ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
				l.IPNet = ip
				l.Iface = &net.Interface{
					Index:        netIface.Index,
					MTU:          netIface.MTU,
					Name:         netIface.Name,
					HardwareAddr: netIface.HardwareAddr,
					Flags:        netIface.Flags,
				}
				l.MAC = netIface.HardwareAddr
				goto end
			}
		}
	}
end:
}

func (l *Local) String() string {
	return fmt.Sprintf("Local(hostname: %s, interface name: %s, MAC addr: %s, IP: %s, Mask: %s)",
		l.Hostname, l.NetIfaceName, l.Iface.HardwareAddr.String(), l.IPNet.IP, l.IPNet.Mask)
}
