package arp

import (
	"net"
	"reflect"
	"testing"
)

func Test_parseIPString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want net.IP
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			args: args{
				s: "192.168.1.3",
			},
			want: []byte{byte(192), byte(168), byte(1), byte(3)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseIPString(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIPString() = %v, want %v", got, tt.want)
			}
		})
	}
}
