module github.com/yeqown/arp

require (
	github.com/google/gopacket v1.1.15
	github.com/mdlayher/raw v0.0.0-20181016155347-fa5ef3332ca9 // indirect

	github.com/yeqown/protocol-impl/ethernet v0.0.0
	golang.org/x/sys v0.0.0-20190124100055-b90733256f2e
)

replace github.com/yeqown/protocol-impl/ethernet => ../ethernet
