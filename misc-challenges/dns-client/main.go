package main

import (
	"fmt"
	"net"
)

type dnsResponse struct {
	header Header
}

type header struct {
	ID               [2]byte
	qrOpcodeAA       [1]byte
	TC_RD_RA_Z_Rcode [1]byte
	qdcount          [2]byte
	ancount          [2]byte
	nscount          [2]byte
	arcount          [2]byte
}

func main() {
	// ips, err := net.LookupAddr("google.com")
	ips, err := net.LookupHost("google.com")
	if err != nil {
		fmt.Print("net.LookupHost err: ", err)
		return
	}
	fmt.Printf("ips: %s\n", ips)

	hosts, err := net.LookupAddr("216.58.195.78")
	if err != nil {
		fmt.Print("net.LookupAddr err: ", err)
		return
	}
	fmt.Printf("hosts: %s\n", hosts)
}
