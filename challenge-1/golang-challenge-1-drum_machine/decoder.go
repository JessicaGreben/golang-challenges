package drum

import (
	"bytes"
	"fmt"
	"os"
)

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
type Pattern struct {
	Header string
	Tempo  string
	Tracks string
}

// String formats the return of the string
// method for the Patter struct.
func (p Pattern) String() string {
	return fmt.Sprintf("%s\n%s",
		p.Header,
		p.Tracks,
	)
}

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.Read failed for file: %s. Error: %v", path, err)
	}
	defer fd.Close()

	fi, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("os.Stat failed: %v", err)
	}

	fileBytes, err := readNextBytes(fd, int(fi.Size()))
	if err != nil {
		return nil, fmt.Errorf("readNextByte failed: %v", err)
	}

	buffer := bytes.NewBuffer(fileBytes)

	drumHeader, err := decodeHeader(buffer)
	if err != nil {
		return nil, fmt.Errorf("decodeHeader failed: %v", err)
	}

	allTracks, err := decodeTracks(buffer)
	if err != nil {
		return nil, fmt.Errorf("decodeTracks failed: %v", err)
	}

	p := Pattern{
		Header: drumHeader,
		Tracks: allTracks,
	}

	return &p, nil
}
