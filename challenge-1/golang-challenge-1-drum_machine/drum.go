// Package drum is implements the decoding of .splice drum machine files.
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

// DrumHeader is the representation of the header of 
// the drum pattern describing the version and tempo.
type DrumHeader struct {
	Version [32]byte
	Tempo float32
}

// String formats the return of the string
// method for the DrumHeader struct.
func (d DrumHeader) String() string {
	return fmt.Sprintf("Saved with HW Version: %s\nTempo: %s",
		bytes.Trim(d.Version[:], "\x00"),
		strings.TrimSuffix(fmt.Sprintf("%.1f", d.Tempo), ".0"),
	)
}

// Header is the representation of the header of 
// the drum pattern describing the format and version.
type Header struct {
	Format  [14]byte
	Version [32]byte
}

// Tempo is the representation of
// the tempo of the drum pattern.
type Tempo struct {
	Tempo float32
}

// TrackHeader is the representation of
// the header for a track pattern.
type TrackHeader struct {
	ID         uint8
	NameLength uint32
}

// Measure is the representation of a measure in
// a drum pattern. Each measure has 16 steps.
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

// String formats the return of the string
// method for the Track struct.
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

// readNextBytes reads the number of bytes from the file.
func readNextBytes(file *os.File, number int) ([]byte, error) {
	bytes := make([]byte, number)
	if _, err := file.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// decodeHeader decodes the header of the track
// all tracks into a single string.
func decodeHeader(buffer *bytes.Buffer) (string, error) {
	var h Header
	if err := binary.Read(buffer, binary.BigEndian, &h); err != nil {
		return "", err
	}

	var t Tempo
	if err := binary.Read(buffer, binary.LittleEndian, &t); err != nil {
		return "", err
	}

	var d = DrumHeader{
		Version: h.Version,
		Tempo: t.Tempo,
	}

	return fmt.Sprint(d), nil
}


// decodeTracks decodes each track and appends
// all tracks into a single string.
func decodeTracks(buffer *bytes.Buffer) (string, error) {
	var allTracks bytes.Buffer
	for {
		var trackHeader TrackHeader
		if err := binary.Read(buffer, binary.BigEndian, &trackHeader); err != nil {
			if err == io.EOF {
				return allTracks.String(), nil
			}
			return "", err
		}

		trackName := make([]byte, trackHeader.NameLength)
		buffer.Read(trackName)

		var measure Measure
		if err := binary.Read(buffer, binary.BigEndian, &measure); err != nil {
			if err == io.EOF {
				return allTracks.String(), nil
			}
			return "", err
		}

		track := Track{
			ID:    trackHeader.ID,
			Name:  string(trackName),
			Steps: fmtBeats(measure.Steps),
		}
		allTracks.WriteString(fmt.Sprint(track))
	}
}

// fmtBeats converts binary representation of the 16 step measure
// pattern into a visualization showing when sound is triggered.
func fmtBeats(steps [16]byte) string {
	var beats bytes.Buffer
	for i := range steps {
		switch steps[i] {
		case 1:

			// "x" represents sound output being triggered in a step.
			beats.WriteString("x")
		default:

			// "-" represents no sound output being triggered in a step.
			beats.WriteString("-")
		}
	}
	return beats.String()
}
