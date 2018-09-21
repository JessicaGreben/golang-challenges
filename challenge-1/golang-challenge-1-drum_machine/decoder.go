package drum

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Pattern is the high level representation of the drum pattern contained
// in a .splice file.
type Pattern struct {
	Header Header
	Tracks []Track
}

// String formats the return of the string method for the Patter struct.
func (p Pattern) String() string {
	var sb strings.Builder
	sb.WriteString(p.Header.String())
	for _, t := range p.Tracks {
		sb.WriteString(t.String())
	}
	return sb.String()
}

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (Pattern, error) {
	fd, err := os.Open(path)
	if err != nil {
		return Pattern{}, fmt.Errorf("os.Read failed for file: %s. Error: %v", path, err)
	}
	defer fd.Close()

	fi, err := os.Stat(path)
	if err != nil {
		return Pattern{}, fmt.Errorf("os.Stat failed: %v", err)
	}

	// Pull in the entire contents of the file.
	data := make([]byte, int(fi.Size()))
	if _, err := fd.Read(data); err != nil {
		return Pattern{}, fmt.Errorf("fd.Read failed: %v", err)
	}

	// Place the raw data into a buffer for processing.
	buffer := bytes.NewBuffer(data)

	// Decode the header section of the data.
	header, err := decodeHeader(buffer)
	if err != nil {
		return Pattern{}, fmt.Errorf("decodeHeader failed: %v", err)
	}

	// Decode the track section of the data.
	tracks, err := decodeTracks(buffer)
	if err != nil {
		return Pattern{}, fmt.Errorf("decodeTracks failed: %v", err)
	}

	p := Pattern{
		Header: header,
		Tracks: tracks,
	}

	return p, nil
}
