package main

import (
	"crypto/rand"
	"fmt"
)

const idLength = 20

type nodeID [idLength]byte

func main() {
	id := generateNodeID()
	fmt.Println(id)
}

func generateNodeID() nodeID {
	id := nodeID{}
	idSlice := make([]byte, idLength)
	rand.Read(idSlice)
	copy(id[:], idSlice)
	return id
}
