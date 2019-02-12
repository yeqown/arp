# ARP Protocol

[![Go Report Card](https://goreportcard.com/badge/github.com/yeqown/arp)](https://goreportcard.com/report/github.com/yeqown/arp) [![GoReportCard](https://godoc.org/github.com/yeqown/arp?status.svg)](https://godoc.org/github.com/yeqown/arp)

The **Address Resolution Protocol (ARP)** is a communication protocol used for discovering the link layer address, such as a MAC address, associated with a given internet layer address, typically an IPv4 address. This mapping is a critical function in the Internet protocol suite. ARP was defined in 1982 by RFC 826,[1] which is Internet Standard STD 37.

ARP has been implemented with many combinations of network and data link layer technologies, such as IPv4, Chaosnet, DECnet and Xerox PARC Universal Packet (PUP) using IEEE 802 standards, FDDI, X.25, Frame Relay and Asynchronous Transfer Mode (ATM). IPv4 over IEEE 802.3 and IEEE 802.11 is the most common usage.

In Internet Protocol Version 6 (IPv6) networks, the functionality of ARP is provided by the Neighbor Discovery Protocol (NDP).

## Internet Protocol (IPv4) over Ethernet ARP packet

|Octet offset|	0|	1|
|-|-------------|----------|
|0|Hardware type (HTYPE)|
|2|Protocol type (PTYPE)|
|4|Hardware address length (HLEN)|	Protocol address length (PLEN)|
|6|Operation (OPER)|
|8|Sender hardware address (SHA) (first 2 bytes)|
|10|(next 2 bytes)|
|12|(last 2 bytes)|
|14|Sender protocol address (SPA) (first 2 bytes)|
|16|(last 2 bytes)|
|18|Target hardware address (THA) (first 2 bytes)|
|20|(next 2 bytes)|
|22|(last 2 bytes)|
|24|Target protocol address (TPA) (first 2 bytes)|
|26|(last 2 bytes)|

## `Socket` in `Go`

### Transport Layer Socket (TCP etc)

We based on the network layer IP protocol and can not customize the socket IP protocol head, called the transport layer socket, it needs to care about the transport layer protocol head how to package, do not need to care about the IP protocol head how to package. It is "theoretically" capable of intercepting any transport layer protocol, and it is also capable of customizing any transport layer protocol, such as a custom protocol called YCP, which is the same level as protocols such as TCP/UDP/ICMP.

ICMP over Transport Layer Socket:
```go
func main() {
    netaddr, _ := net.ResolveIPAddr("ip4", "172.17.0.3")
    conn, _ := net.ListenIP("ip4:icmp", netaddr)
    for {
        buf := make([]byte, 1024)
        n, addr, _ := conn.ReadFrom(buf)
        msg,_:=icmp.ParseMessage(1,buf[0:n])
        fmt.Println(n, addr, msg.Type,msg.Code,msg.Checksum)
    }
}
```

### Network Layer Socket(IP)

We based on the network layer IP protocol and can customize the socket IP protocol head, known as the network layer socket, it needs to care about the IP protocol head how to package, do not need to care about the Ethernet frame head and tail how to package

with go stdlib:

```go
func main() {
    netaddr, _ := net.ResolveIPAddr("ip4", "172.17.0.3")
    conn, _ := net.ListenIP("ip4:tcp", netaddr)

    // transfer transport layer socket into raw socket
    ipconn,_:=ipv4.NewRawConn(conn)
    for {
        buf := make([]byte, 1480)
        hdr, payload, controlMessage, _ := ipconn.ReadFrom(buf)
        fmt.Println("ipheader:",hdr,controlMessage)
        tcpheader:=NewTCPHeader(payload)
        fmt.Println("tcpheader:",tcpheader)
    }
}
```

without go stdlib, `syscall` instead:
```go
func main() {
    fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
    f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
    for {
        buf := make([]byte, 1500)
        f.Read(buf)
        ip4header, _ := ipv4.ParseHeader(buf[:20])
        fmt.Println("ipheader:", ip4header)
        tcpheader := util.NewTCPHeader(buf[20:40])
        fmt.Println("tcpheader:", tcpheader)
    }
}
```

### Link Layer Socket (Ethernet)


```go
func main() {
    ifi, err := net.InterfaceByName("eth0")
    util.CheckError(err)

    // 0x800 means Internet Protocol packet
    // syscall.Socket(domain, typ, proto int)
    fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(v2.Htons(0x800)))
    for {
        buf := make([]byte, 1514)
        n, _, _ := syscall.Recvfrom(fd, buf, 0)
        header := new(v2.Header)
        header.Decode(buf[0:14])
        fmt.Println(header)
    }  
}
```

## Examples

with lib `gopacket`:
```go
```

with socket:
```go
```

## Reference

* [Go中原始套接字的深度实践](https://www.cnblogs.com/mushroom/p/9097409.html)
* [Go中链路层套接字的实践](https://www.cnblogs.com/mushroom/p/9321190.html)