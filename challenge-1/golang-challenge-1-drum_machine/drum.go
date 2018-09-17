// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Header of the drum pattern
type Header struct {
	_       [14]byte
	Version [32]byte
}

func (h Header) String() string {
	// convert to byte type in order to trim padding
	version := h.Version[:]
	return fmt.Sprintf("Saved with HW Version: %s\n",
		bytes.Trim(version, "\x00"),
	)
}

// Tempo of the drum pattern
type Tempo struct {
	Tempo float32
}

func (t Tempo) String() string {
	tempo := strings.TrimSuffix(fmt.Sprintf("%.1f", t.Tempo), ".0")
	return fmt.Sprintf("Tempo: %s\n", tempo)
}

// TrackHeader is the header for a track.
type TrackHeader struct {
	ID         uint8
	NameLength uint32
}

// Measure of a track.
// Each measure has 16 steps.
type Measure struct {
	Steps [16]byte
}

// Track with the ID, Name, and formatted steps.
type Track struct {
	ID    uint8
	Name  string
	Steps string
}

func (t Track) String() string {
	return fmt.Sprintf("(%d) %s\t|%s|%s|%s|%s|\n",
		t.ID,
		string(t.Name),
		t.Steps[0:4],
		t.Steps[4:8],
		t.Steps[8:12],
		t.Steps[12:16],
	)
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func decodeHeader(buffer *bytes.Buffer) (string, string) {
	h := Header{}
	err := binary.Read(buffer, binary.BigEndian, &h)
	if err != nil {
		fmt.Print("Header: binary.Read failed ", err)
	}

	t := Tempo{}
	err = binary.Read(buffer, binary.LittleEndian, &t)
	if err != nil {
		fmt.Print("Tempo: binary.Read failed ", err)
	}
	return fmt.Sprint(h), fmt.Sprint(t)
}

func decodeTracks(buffer *bytes.Buffer) string {
	allTracks := ""
	for {
		trackHeader := TrackHeader{}
		err := binary.Read(buffer, binary.BigEndian, &trackHeader)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Print("Tracks: binary.Read failed:", err)
		}

		trackName := make([]byte, trackHeader.NameLength)
		buffer.Read(trackName)

		measure := Measure{}
		err = binary.Read(buffer, binary.BigEndian, &measure)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Print("Steps: binary.Read failed:", err)
		}

		beats := fmtBeats(measure.Steps)
		track := Track{trackHeader.ID, string(trackName), beats}
		allTracks = fmt.Sprintf(`%s%s`, allTracks, track)
	}
	return allTracks
}

// fmtBeats converts binary representation of the 16 step measure
// pattern into a visualization showing when sound is triggered
func fmtBeats(steps [16]byte) string {
	beats := ""
	for i := 0; i < 16; i++ {
		if steps[i] == 1 {
			// "x" represents sound output being triggered in a step
			beats += "x"
		} else {
			// "-" represents no sound output being triggered in a step
			beats += "-"
		}
	}
	return beats
}
