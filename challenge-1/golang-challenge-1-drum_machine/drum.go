// Package drum is supposed to implement the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

// Header is the representation of
// the header of the drum pattern
// describing the format and version
type Header struct {
	Format  [14]byte
	Version [32]byte
}

func (h Header) String() string {
	return fmt.Sprintf("Saved with HW Version: %s",
		bytes.Trim(h.Version[:], "\x00"),
	)
}

// Tempo is the representation of
// the tempo of the drum pattern
type Tempo struct {
	Tempo float32
}

func (t Tempo) String() string {
	tempo := strings.TrimSuffix(fmt.Sprintf("%.1f", t.Tempo), ".0")
	return fmt.Sprintf("Tempo: %s", tempo)
}

// TrackHeader is the representation of
// the header for a track pattern.
type TrackHeader struct {
	ID         uint8
	NameLength uint32
}

// Measure is the representation of
// a measure in a drum pattern.
// Each measure has 16 steps.
type Measure struct {
	Steps [16]byte
}

// Track is the string representation of a track
// combining the ID, Name, and formatted steps.
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

func readNextBytes(file *os.File, number int) ([]byte, error) {
	bytes := make([]byte, number)
	if _, err := file.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

func decodeHeader(buffer *bytes.Buffer) (string, string, error) {
	var h = Header{}
	if err := binary.Read(buffer, binary.BigEndian, &h); err != nil {
		return "", "", err
	}

	var t = Tempo{}
	if err := binary.Read(buffer, binary.LittleEndian, &t); err != nil {
		return "", "", err
	}

	return fmt.Sprint(h), fmt.Sprint(t), nil
}

func decodeTracks(buffer *bytes.Buffer) (string, error) {
	var allTracks = ""
	for {
		trackHeader := TrackHeader{}
		if err := binary.Read(buffer, binary.BigEndian, &trackHeader); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		trackName := make([]byte, trackHeader.NameLength)
		buffer.Read(trackName)

		measure := Measure{}
		if err := binary.Read(buffer, binary.BigEndian, &measure); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		beats := fmtBeats(measure.Steps)
		track := Track{
			ID:    trackHeader.ID,
			Name:  string(trackName),
			Steps: beats,
		}
		allTracks = fmt.Sprintf(`%s%s`, allTracks, track)
	}

	return allTracks, nil
}

// fmtBeats converts binary representation of the 16 step measure
// pattern into a visualization showing when sound is triggered
func fmtBeats(steps [16]byte) string {
	var beats = ""
	for i := range steps {
		switch steps[i] {
		case 1:

			// "x" represents sound output being triggered in a step
			beats += "x"
		default:

			// "-" represents no sound output being triggered in a step
			beats += "-"
		}
	}
	return beats
}
