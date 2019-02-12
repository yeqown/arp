package v2

import (
	"errors"
	"fmt"
	"net"
	"syscall"
	"time"
	
	"golang.org/x/net/ipv4"
	"github.com/yeqown/protocol-impl/ethernet"
)

var (
	errEmptyInterface = errors.New("empty interface")
)

// ARPCall ...
func ARPCall(ifi *net.Interface) error {
	if ifi == nil {
		return errEmptyInterface
	}

	fd, err := syscall.Socket(syscall.AF_APPLETALK, syscall.SOCK_RAW, syscall.IPPROTO_ETHERIP)
	if err != nil {
		return fmt.Errorf("syscall.Socket %v", err)
	}
	data := []byte("msg")
	minPayload := len(data)
	if minPayload < 46 {
		minPayload = 46
	}
	b := make([]byte, 14+minPayload)
	broadcast := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	header := &ethernet.Header{
		DestinationAddress: broadcast,
		SourceAddress:      ifi.HardwareAddr,
		EthernetType:       syscall.IPPROTO_ETHERIP,
	}
	copy(b[0:14], header.Encode())
	copy(b[14:14+len(data)], data)

	var baddr [8]byte
	copy(baddr[:], broadcast)
	to := &sys.SockaddrLinklayer{
		Ifindex:  ifi.Index,
		Halen:    6,
		Addr:     baddr,
		Protocol: syscall.IPPROTO_ETHERIP,
	}
	for {
		log.Printf("sending %s\n", b)
		err = syscall.Sendto(fd, b, 0, to)
		if err != nil {
			return fmt.Errorf("syscall.Sendto got err: %v", err)
		}
		time.Sleep(time.Second)
	}
	return nil
}
