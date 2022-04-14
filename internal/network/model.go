package network

import "net"

type NetworkAddress struct {
	Interface net.Interface
	IP        net.IP
	Net       *net.IPNet
}
