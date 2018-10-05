package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/nacl/box"
)

// SecureReader implements the io.Reader interface to decrypt encrypted bytes.
type SecureReader struct {
	io.Reader
	privateKey [32]byte
	publicKey  [32]byte
}

// NewSecureReader is a factory function for SecureReader.
func NewSecureReader(r io.Reader, privateKey, publicKey [32]byte) *SecureReader {
	return &SecureReader{
		Reader:     r,
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// Read implements the io.Reader interface for secureReader to decrypt bytes.
func (sr *SecureReader) Read(p []byte) (int, error) {
	n, err := decrypt(p, sr)
	if err != nil {
		return 0, err
	}

	return n, nil
}

// decrypt decrypts bytes using public-key cryptography.
func decrypt(p []byte, sr *SecureReader) (int, error) {

	// Read the nonce.
	var nonce [24]byte
	n, err := io.ReadFull(sr.Reader, nonce[:])
	if err != nil {
		return 0, err
	}

	// Overhead is the number of bytes of overhead when
	// encrypting a message therefore we need to create a
	// buffer larger than just the encrypted message.
	msg := make([]byte, len(p)+box.Overhead)

	// Read the encrypted message.
	n, err = sr.Reader.Read(msg)
	if err != nil {
		return 0, err
	}

	// The len and capacity of p will be set by the caller. On the call to
	// box.Open the length value of p must be set to 0 since box. Open will
	// append values to the slice. So by using p[:0] we are passing a slice
	// value who's length is 0 but capacity represents the full size of the
	// backing array. This will allow box.Open to use append to add the data
	// during the decryption.

	// box.Open will decypt the message and place just the decypted bytes back
	// to the p slice. It will also return just those decypted bytes in the
	// returned slice. We need the returned slice to get the length of bytes
	// for this message. Because p will still have a larger length than what
	// was appended from index 0.
	dec, ok := box.Open(p[:0], msg[:n], &nonce, &sr.publicKey, &sr.privateKey)
	if !ok {
		return 0, err
	}

	return len(dec), nil
}

// SecureWriter implements the io.Writer interface to encrypt bytes.
type SecureWriter struct {
	io.Writer
	privateKey [32]byte
	publicKey  [32]byte
}

// NewSecureWriter is a factory function for SecureWriter.
func NewSecureWriter(w io.Writer, privateKey, publicKey [32]byte) *SecureWriter {
	return &SecureWriter{
		Writer:     w,
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// Write implements the io.Writer interface for secureWriter to encrypt bytes.
func (sw *SecureWriter) Write(p []byte) (int, error) {
	encryptedMsg, err := encrypt(p, sw.publicKey, sw.privateKey)
	if err != nil {
		return 0, err
	}

	// Write wants to know that we processed all the bytes in p,
	// so we need to report that we did with len(p)
	sw.Writer.Write(encryptedMsg)
	return len(p), nil
}

// encrypt encrypts bytes using public-key cryptography.
func encrypt(p []byte, publicKey, privateKey [32]byte) ([]byte, error) {

	// Create a nonce to encrypt with.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	// Encrypt the message which appends the result to the nonce
	// in order to store the nonce with the encrypted message so that we
	// can use the same nonce when its decrypted.
	encryptedMsg := box.Seal(nonce[:], p, &nonce, &publicKey, &privateKey)
	return encryptedMsg, nil
}

// Perform a pubic key exchange over the connection.
func handshake(conn net.Conn, senderKey [32]byte) ([32]byte, error) {
	var receiverKey [32]byte

	// Send the public key.
	if _, err := conn.Write(senderKey[:]); err != nil {
		return receiverKey, err
	}

	// Recieve the public key from the other side
	// of the connection.
	if _, err := conn.Read(receiverKey[:]); err != nil {
		return receiverKey, err
	}

	return receiverKey, nil
}

// SecureConn provides secure read and write over a connection.
type SecureConn struct {
	net.Conn
	*SecureReader // Read method is promoting up.
	*SecureWriter // Write method is promoting up.
}

// Read adds a secure read that decrypts messages before reading.
func (sc *SecureConn) Read(p []byte) (int, error) {
	return sc.SecureReader.Read(p)
}

// Write adds a secure write that encrypts messages before writing.
func (sc *SecureConn) Write(p []byte) (int, error) {
	return sc.SecureWriter.Write(p)
}

// Dial connects to the server over a secure connections which encrypts and
// decrypts all read and write bytes.
func Dial(addr string) (net.Conn, error) {

	// Connect to server.
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Generate a public/private encryption keypair for the client.
	// GenerateKey is using pointer semantics for the keys. We want to
	// use value sematics. We think it is the right semantic for this data.
	clientPublicKey, clientPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Perform a public encryption key exchange to
	// get the servers' public encryption key.
	serverPublicKey, err := handshake(conn, *clientPublicKey)
	if err != nil {
		return nil, err
	}

	// Wrap the connection with a secure reader and secure writer
	// which will encrypt/decrypt all read/writes.
	sc := SecureConn{
		Conn:         conn,
		SecureReader: NewSecureReader(conn, *clientPrivateKey, serverPublicKey),
		SecureWriter: NewSecureWriter(conn, *clientPrivateKey, serverPublicKey),
	}

	return &sc, nil
}

// Serve starts a secure echo server on the given listener.
func Serve(l net.Listener) error {

	// Accept a connection from a client.
	conn, err := l.Accept()
	if err != nil {
		return fmt.Errorf("Error l.Accept : %v", err)
	}
	defer conn.Close()

	// Generate a public/private encryption keypair for the server.
	serverPublicKey, serverPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("Error Serve GenKey: %v", err)
	}

	// Perform a public encryption key exchange to
	// get the client's public encryption key.
	clientPubKey, err := handshake(conn, *serverPublicKey)
	if err != nil {
		return fmt.Errorf("Error Serve handshake: %v", err)
	}

	// Construct the secure reader and writer. We have to convert the keys
	// from pointer to value semantics for our code base.
	sw := NewSecureWriter(conn, *serverPrivateKey, clientPubKey)
	sr := NewSecureReader(conn, *serverPrivateKey, clientPubKey)

	// Echo the data back to the client over a secure
	// connection that encrypts and decrypts all reads/writes.
	if _, err = io.Copy(sw, sr); err != nil {
		return fmt.Errorf("Error io.Copy: %v", err)
	}

	return nil
}

func main() {
	port := flag.Int("l", 0, "Listen mode. Specify port")
	flag.Parse()

	// Server mode.
	if *port != 0 {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			fmt.Printf("Error net.Listen: %v", err)
		}
		defer l.Close()
		log.Fatal(Serve(l))
	}

	// Client mode.
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <port> <message>\n", os.Args[0])
	}
	conn, err := Dial("localhost:" + os.Args[1])
	if err != nil {
		fmt.Printf("Error Dial: %v", err)
	}
	if _, err := conn.Write([]byte(os.Args[2])); err != nil {
		fmt.Printf("Error conn.Write: %v", err)
	}
	buf := make([]byte, len(os.Args[2]))
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("Error conn.Read: %v", err)
	}
	fmt.Printf("%s\n", buf[:n])
}
