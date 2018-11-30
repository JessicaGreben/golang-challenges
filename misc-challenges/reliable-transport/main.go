package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
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

	addr := net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: 1337,
	}

	// Listen for incoming UDP packets.
	conn, err := net.ListenUDP("udp", &addr)
	// conn, err := net.ListenPacket("udp", "localhost:1337")
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Print("Listening on port 1337...\n")

	for {
		buffer := make([]byte, 4096)
		// n, addr, err := conn.ReadFrom(buffer)
		// n, addr, err := conn.ReadFromUDP(buffer)
		oob := make([]byte, 1024)
		n, on, _, addr, err := conn.ReadMsgUDP(buffer, oob)
		if err != nil {
			// if io.EOF
			return err
		}

		// Process UDP packet.
		fmt.Printf("length of oob: %d\n", on)
		fmt.Printf("string: %s\n", buffer[:n])
		ID := buffer[:9]
		// ID, err := process(buffer[:n])
		// if err != nil {
		// 	fmt.Print(err)
		// }

		// ACK response.
		resp := fmt.Sprintf("ACK: %d", ID)
		conn.WriteTo([]byte(resp), addr)
	}
}

func process(packet []byte) (uint64, error) {
	fmt.Printf("Payload received: %v\n", packet)

	r := bytes.NewReader(packet)
	p := packetRUDP{}
	err := binary.Read(r, binary.BigEndian, &p)
	fmt.Printf("String: %#v\n", p)

	if err != nil {
		return 0, err
	}
	return p.ID, nil
}

func client(port int) error {

	addr := net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: port,
	}

	// Connect UDP.
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	packet := fmt.Sprintf("%d: Packet payload.", rand.Uint64())
	fmt.Printf("Sending packet: %s\n", packet)
	_, err = conn.Write([]byte(packet))
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	n, _, _, _, err := conn.ReadMsgUDP(buffer, nil)
	if err != nil {
		return err
	}
	fmt.Printf("Read ACK: %s\n", buffer[:n])

	return nil
}

// func waitingForACK(conn *net.UDPConn) bool {
// 	// TODO: start timer to listen for an ACK.
// 	t := time.NewTimer(2 * time.Second)
// 	buffer := make([]byte, 1024)
// 	n, _, _, _, _ := conn.ReadMsgUDP(buffer, nil)
// 	// conn.Read(buffer[:n])
// 	// Resend if timeout is reached.
// 	tStop := t.Stop()
// 	if tStop {
// 		return false
// 	}
// 	return true
// }

// func timeout() {
// 	timer1 := time.NewTimer(2 * time.Second)
// 	go func() {
// 		<-timer1.C
// 		fmt.Println("Timer expired: No ACK")
// 	}()
// }

type packetRUDP struct {
	ID   uint64
	Data [1]byte
}

type headersRUDP struct {
	ID     uint64
	Length uint16
}

// func (p packetRUDP) send(conn *net.UDPConn) {
// 	var arr [1]byte
// 	copy(arr[:], "abc")
// 	p.Data = arr
// 	p.addHeaders()

// 	enc := gob.NewEncoder(conn)
// 	err := enc.Encode(p)
// 	if err != nil {
// 		fmt.Println("encode error:", err)
// 		return
// 	}

// buf := new(bytes.Buffer)
// err := binary.Write(buf, binary.BigEndian, &p)

// if err != nil {
// 	fmt.Println("binary.Write failed:", err)
// 	return
// }
// err := binary.Read(r, binary.BigEndian, &p)

// conn.WriteMsgUDP(, nil, nil)
// }

// func (p packetRUDP) addHeaders() {
// 	p.ID = rand.Uint64()
// 	p.Length = uint16(len(p.Data))
// 	fmt.Printf("ID: %d\n", p.ID)
// 	fmt.Printf("length of data is: %d\n", p.Length)
// }
