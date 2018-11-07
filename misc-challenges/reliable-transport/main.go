package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Specify -c for client or -s for server.")
		fmt.Println("Example: go run main.go -c <port>, where port is the port of the unreliable proxy server.")
		fmt.Println("Example: go run main.go -s. The server listens on port 1337 so no need to provide a port value.")
		return
	}

	switch os.Args[1] {
	case "-s":
		server()
	case "-c":
		port, _ := strconv.Atoi(os.Args[2])
		client(port)
	}
}

func server() error {

	// Listen for incoming UDP packets.
	conn, err := net.ListenPacket("udp", "localhost:1337")
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Print("Listening on port 1337...\n")

	for {
		buffer := make([]byte, 4096)
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			return err
		}

		// Process UDP packet.
		process(buffer[:n])

		// ACK response.
		conn.WriteTo([]byte("ACK from server."), addr)
	}
}

func process(packetUDP []byte) {
	fmt.Printf("Payload received: %s", packetUDP)
	// read header
	// headers, payload := readHeadersUDP(packetUDP)
	// checksum header
}

func client(port int) error {

	n := net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: port,
	}

	// Connect UDP.
	conn, err := net.DialUDP("udp", nil, &n)
	if err != nil {
		return err
	}
	defer conn.Close()

	sendPacketUDP(conn)

	waitForACK(conn, n)
	// buffer := make([]byte, 1024)
	// conn.Read(buffer)
	// fmt.Printf("buf: %v\n", buffer)
	timeout()

	return nil
}

func sendPacketUDP(conn *net.UDPConn) {
	payload := []byte("Payload from client.\n")
	msg := addHeadersUDP(payload)
	conn.WriteMsgUDP(msg, nil, nil)
}

func addHeadersUDP(payload []byte) []byte {
	return payload
}

func waitForACK(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	n, _, _ := conn.ReadMsgUDP(buffer)
	conn.Read
	// TODO: start timer to listen for an ACK.
	// Resend if timeout is reached.

}

func timeout() {
	timer1 := time.NewTimer(2 * time.Second)
	go func() {
		<-timer1.C
		fmt.Println("Timer expired: No ACK")
	}()
}

type packetRUDP struct {
	headers [2]byte
	data    [1]byte
}

type headersRUDP struct{}
