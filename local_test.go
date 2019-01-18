// Package arp ...
package arp

import (
	"net"
	"reflect"
	"testing"
)

func xTestNewLocal(t *testing.T) {
	iface, _ := net.InterfaceByName("en0")

	type args struct {
		netIfaceName string
	}
	tests := []struct {
		name string
		args args
		want *Local
	}{
		// TODO: Add test cases.
		{
			name: "case 0",
			args: args{
				netIfaceName: "en0",
			},
			want: &Local{
				NetIfaceName: "en0",
				Hostname:     "192.168.0.130",
				MAC:          iface.HardwareAddr,
				IPNet:        &net.IPNet{},
				Iface:        iface,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLocal(tt.args.netIfaceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}
