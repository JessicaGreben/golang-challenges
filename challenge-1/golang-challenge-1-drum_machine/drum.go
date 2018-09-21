// Package drum is implements the decoding of .splice drum machine files.
// See golang-challenge.com/go-challenge1/ for more information
package drum

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// Header is the representation of the header of the drum pattern
// describing the version and tempo.
type Header struct {
	Version [32]byte
	Tempo   float32
}

// String formats the return of the string method for the Header struct.
func (d Header) String() string {
	return fmt.Sprintf("Saved with HW Version: %s\nTempo: %s\n",
		bytes.Trim(d.Version[:], "\x00"),
		strings.TrimSuffix(fmt.Sprintf("%.1f", d.Tempo), ".0"),
	)
}

// decodeHeader decodes the header of the drum pattern into a Header struct.
func decodeHeader(buffer *bytes.Buffer) (Header, error) {

	// Extract the header bytes from the data.
	var header struct {
		Format  [14]byte
		Version [32]byte
	}
	if err := binary.Read(buffer, binary.BigEndian, &header); err != nil {
		return Header{}, err
	}

	// Extract the tempo value which is the next four bytes.
	var tempo float32
	if err := binary.Read(buffer, binary.LittleEndian, &tempo); err != nil {
		return Header{}, err
	}

	h := Header{
		Version: header.Version,
		Tempo:   tempo,
	}

	return h, nil
}

// Track is the string representation of a track combining the ID, Name,
// and formatted steps.
type Track struct {
	ID    uint8
	Name  string
	Steps [16]byte
}

// String formats the return of the string method for the Track struct.
func (t Track) String() string {

	// Convert the bytes representation of the steps
	// into the desired string format.
	var sb strings.Builder
	for i := range t.Steps {
		switch t.Steps[i] {
		case 1:

			// "x" represents sound output being triggered in a step.
			sb.WriteString("x")
		default:

			// "-" represents no sound output being triggered in a step.
			sb.WriteString("-")
		}
	}
	steps := sb.String()

	return fmt.Sprintf("(%d) %s\t|%s|%s|%s|%s|\n",
		t.ID,
		t.Name,
		steps[0:4],
		steps[4:8],
		steps[8:12],
		steps[12:16],
	)
}

// decodeTracks decodes each track and appends all tracks into a single slice.
func decodeTracks(buffer *bytes.Buffer) ([]Track, error) {
	var tracks []Track

	for {

		// Extract the header of the track.
		var header struct {
			ID     uint8
			Length uint32
		}
		if err := binary.Read(buffer, binary.BigEndian, &header); err != nil {
			if err == io.EOF {
				return tracks, nil
			}
			return nil, err
		}

		// Usee the value of the header.length to extract
		// the name of the track.
		name := make([]byte, header.Length)
		buffer.Read(name)

		// Extract the measure steps which are the next 16 bytes.
		var steps [16]byte
		if err := binary.Read(buffer, binary.BigEndian, &steps); err != nil {
			if err == io.EOF {
				return tracks, nil
			}
			return nil, err
		}

		track := Track{
			ID:    header.ID,
			Name:  string(name),
			Steps: steps,
		}
		tracks = append(tracks, track)
	}
}
