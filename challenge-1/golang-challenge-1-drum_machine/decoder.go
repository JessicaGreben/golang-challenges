package drum

import (
	"bytes"
	"fmt"
	"os"
)

type Pattern struct {
	Header string
	Tempo  string
	Tracks string
}

func (p Pattern) String() string {
	return fmt.Sprintf("%s%s%s",
		p.Header,
		p.Tempo,
		p.Tracks,
	)
}

// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
func DecodeFile(path string) (*Pattern, error) {
	fd, err := os.Open(path)
	if err != nil {
		fmt.Printf("os.Read failed for file: %s. Error: %v", path, err)
	}
	defer fd.Close()

	fi, err := os.Stat(path)
	if err != nil {
		fmt.Print("os.Stat failed: ", err)
	}
	size := fi.Size()
	fileBytes := readNextBytes(fd, int(size))
	buffer := bytes.NewBuffer(fileBytes)

	header, tempo := decodeHeader(buffer)
	allTracks := decodeTracks(buffer)

	p := Pattern{header, tempo, allTracks}
	return &p, nil
}
