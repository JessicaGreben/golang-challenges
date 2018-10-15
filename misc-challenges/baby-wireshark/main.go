package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
)

// The pcap file global header.
// ref: https://www.tcpdump.org/manpages/pcap-savefile.5.txt
type globalHeader struct {
	MagicNumber     [4]byte
	MajorVersion    int16
	MinorVersion    [2]byte
	ZoneOffset      [4]byte
	TimeAccuracy    [4]byte
	SnapshotLength  uint32
	LinkLayerHeader [4]byte
}

// Pcap per-packet header.
type packetHeader struct {
	TimestampSec       [4]byte
	TimestampMs        [4]byte
	PacketSize         uint32
	OriginalPacketSize uint32
}

// Per packet link layer frame header. Layer 2 Ethernet frame.
// ref: https://en.wikipedia.org/wiki/Ethernet_frame#Structure
type frameHeader struct {
	DestAddr   [6]byte
	SourceAddr [6]byte
	EtherType  uint16
}

// Per packet IP layer datagram header.
// ref: https://tools.ietf.org/html/rfc791#page-11
type datagramHeader struct {
	Version        uint8
	_              [1]byte
	TotalLength    uint16
	ID             [2]byte
	_              [2]byte
	TTL            [1]byte
	Protocol       uint8
	HeaderChecksum [2]byte
	Source         [4]byte
	Dest           [4]byte
}

// Per packet transport layer segment header.
// ref: https://tools.ietf.org/html/rfc793#section-3.1
type segmentHeader struct {
	SourcePort uint16
	DestPort   uint16
	Sequence   uint32
	AckNumber  [4]byte
	DataOffset [1]byte
}

func main() {

	fd, err := os.Open("./net.cap")
	if err != nil {
		fmt.Print(err)
	}
	defer fd.Close()

	fi, err := os.Stat("./net.cap")
	if err != nil {
		fmt.Print(err)
	}

	// Read the entire contents of the file.
	data := make([]byte, int(fi.Size()))
	if _, err := fd.Read(data); err != nil {
		fmt.Print(err)
	}

	// Place raw bytes into the buffer for processing.
	buffer := bytes.NewBuffer(data)

	// Read the pcap-savefile global header.
	if err := readGlobalHeader(buffer); err != nil {
		fmt.Print(err)
	}

	// This mapping is the tcp sequence number to http payload bytes
	// so that we can reassemble the bytes in the correct order.
	httpData := make(map[int][]byte)

	// Read all of the packets.
	for {
		httpData, err = readPacket(buffer, httpData)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Print(err)
			return
		}
	}

	// Create a file from the httpData and open that file.
	createFile(httpData)
}

func readPacket(buffer *bytes.Buffer, httpOrder map[int][]byte) (map[int][]byte, error) {
	pacHeader, err := readPacHeader(buffer)
	if err != nil {
		return httpOrder, err
	}

	readFrame(buffer)

	datagram, err := readDatagram(buffer)
	if err != nil {
		return httpOrder, err
	}

	datagramHeader := 4 * int(datagram.Version&0x0f)

	segment, err := readSegment(buffer)
	if err != nil {
		return httpOrder, err
	}
	tcpHeader := 4 * int(segment.DataOffset[0]>>4)

	// Read HTTP header and data.
	// The HTTP payload length is the datagram total length minus the datagram and
	//  tcp header lengths.
	httpLength := int(datagram.TotalLength) - datagramHeader - tcpHeader
	httpData := make([]byte, httpLength)
	buffer.Read(httpData)

	// We know we only want http data from packets
	// that are destined to this addr.
	ip := [4]byte{192, 168, 0, 101}
	if httpLength > 0 && datagram.Dest == ip {

		// The http header is separated by \r\n\r\n. We can tell this
		// packet contains the header if this is present.
		headerPresent := bytes.Split(httpData, []byte{'\r', '\n', '\r', '\n'})
		if len(headerPresent) > 1 {
			httpOrder[int(segment.Sequence)] = headerPresent[1]
		} else {
			httpOrder[int(segment.Sequence)] = httpData
		}
	}

	// Weird hack to remove the random extra 2 bytes that show up randomly
	// on packets of this size.
	if pacHeader.PacketSize == 68 || pacHeader.PacketSize == 76 {
		buffer.Read(make([]byte, 2))
	}

	return httpOrder, nil
}

func readGlobalHeader(buffer *bytes.Buffer) error {
	globalHeader := globalHeader{}
	if err := binary.Read(buffer, binary.BigEndian, &globalHeader); err != nil {
		return err
	}

	// Make sure link layer header is Ethernet.
	// ref: https://www.tcpdump.org/linktypes.html
	if globalHeader.LinkLayerHeader[0] != 1 {
		return fmt.Errorf("Expected link layer header to be Ethernet, but got %d instead.", globalHeader.LinkLayerHeader)
	}

	return nil
}

func readPacHeader(buffer *bytes.Buffer) (packetHeader, error) {
	pacHeader := packetHeader{}
	if err := binary.Read(buffer, binary.LittleEndian, &pacHeader); err != nil {
		if err == io.EOF {
			return pacHeader, err
		}
		fmt.Print("binary.Read err: ", err)
	}

	// This won't always be the case, but for this file we know that the packet
	// size is the same as the original size. I.e. the packet is never trucated.
	if int(pacHeader.PacketSize) != int(pacHeader.OriginalPacketSize) {
		return pacHeader, fmt.Errorf("Error! pcap pac header is wrong. Expected %d, but got %d.", pacHeader.PacketSize, pacHeader.OriginalPacketSize)
	}
	return pacHeader, nil
}

func readFrame(buffer *bytes.Buffer) error {
	frame := frameHeader{}
	if err := binary.Read(buffer, binary.BigEndian, &frame); err != nil {
		return err
	}

	// We know that the ether type is 8, indicating protocol IPv4 for the payload.
	if frame.EtherType != 8 {
		return fmt.Errorf("Error! EtherType is incorrect. Should be 8, but got %d\n", frame.EtherType)
	}
	return nil
}

func readDatagram(buffer *bytes.Buffer) (datagramHeader, error) {
	datagram := datagramHeader{}
	if err := binary.Read(buffer, binary.BigEndian, &datagram); err != nil {
		return datagram, err
	}

	// Make sure the protocol is 6 for TCP.
	if datagram.Protocol != 6 {
		return datagram, fmt.Errorf("IP Protocol should be 6, but it is actially %d\n", datagram.Protocol)
	}

	return datagram, nil
}

func readSegment(buffer *bytes.Buffer) (segmentHeader, error) {
	segment := segmentHeader{}
	if err := binary.Read(buffer, binary.BigEndian, &segment); err != nil {
		return segment, err
	}

	tcpHeader := 4 * int(segment.DataOffset[0]>>4)

	// We only read in 13 bytes so far, but the segment header is bigger then that,
	// so we need to read the rest of it.
	readTheRestTcp := tcpHeader - 13
	buffer.Read(make([]byte, readTheRestTcp))

	return segment, nil
}

func createFile(httpOrder map[int][]byte) error {

	// Sort the values of the TCP sequence numbers.
	var sequenceNums []int
	for k := range httpOrder {
		sequenceNums = append(sequenceNums, k)
	}
	sort.Ints(sequenceNums)

	f, err := os.Create("./packet.jpeg")
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the HTTP data in order to a jpeg.
	for _, k := range sequenceNums {
		data := httpOrder[k]
		f.Write(data)
	}
	f.Sync()

	// Open the jpeg that was created.
	if err := exec.Command("/usr/bin/open", "./packet.jpeg").Run(); err != nil {
		return err
	}

	return nil
}
