package network

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Network interface {
	GetLocalAddresses() ([]NetworkAddress, error)
	GetAllNetworkIps(*NetworkAddress) *[]net.IP
	GetHostName() string
}

func GetLocalAddresses() (*[]NetworkAddress, error) {
	ifaces, err := net.Interfaces()
	var result []NetworkAddress
	if err != nil {
		log.Print(fmt.Errorf("localAddresses: %v", err.Error()))
		return &result, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Print(fmt.Errorf("localAddresses: %v", err.Error()))
			continue
		}
		if strings.Contains(i.Flags.String(), "up") {
			for _, a := range addrs {
				ip, net, _ := net.ParseCIDR(a.String())
				if !net.IP.IsLoopback() && ip.IsPrivate() && net.IP.To4() != nil {
					result = append(result, NetworkAddress{
						Interface: i,
						IP:        ip,
						Net:       net,
					})
				}
			}
		}

	}
	return &result, nil

}

func GetAllNetworkIps(in *NetworkAddress) *[]net.IP {
	var result []net.IP
	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(in.Net.Mask)
	start := binary.BigEndian.Uint32(in.Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	// loop through addresses as uint32
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		result = append(result, ip)
	}
	return &result
}

func GetHostName() string {
	name, err := os.Hostname()
	if err != nil {
		log.Println("Error getting hostname. ", err.Error())
	}
	return name
}
