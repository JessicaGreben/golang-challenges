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

// secureReader implements the io.Reader interface to decrypt encrypted bytes.
type secureReader struct {
	r    io.Reader
	priv *[32]byte
	pub  *[32]byte
}

// NewSecureReader is a factory function for secureReader.
func NewSecureReader(r io.Reader, priv, pub *[32]byte) io.Reader {
	sr := secureReader{
		r:    r,
		priv: priv,
		pub:  pub,
	}
	return &sr
}

// Read implements the io.Reader interface for secureReader to decrypt bytes.
func (sr *secureReader) Read(p []byte) (int, error) {
	dec, err := decrypt(p, sr)
	if err != nil {
		return 0, err
	}

	return len(dec), nil
}

// secureWriter implements the io.Writer interface to encrypt bytes.
type secureWriter struct {
	w    io.Writer
	priv *[32]byte
	pub  *[32]byte
}

// NewSecureWriter is a factory function for secureWriter.
func NewSecureWriter(w io.Writer, priv, pub *[32]byte) io.Writer {
	sw := secureWriter{
		w:    w,
		priv: priv,
		pub:  pub,
	}
	return &sw
}

// Write implements the io.Writer interface for secureWriter to encrypt bytes.
func (sw *secureWriter) Write(p []byte) (int, error) {
	enc, err := encrypt(p, sw.pub, sw.priv)
	if err != nil {
		return 0, err
	}

	sw.w.Write(enc)
	return len(p), nil
}

// secureReadWriteCloser implements the ReadWriteCloser interface.
type secureReadWriteCloser struct {
	io.Reader
	io.Writer
	io.Closer
}

// newSecureReadWriteCloser is a factory function for secureReadWriteCloser.
func newSecureReadWriteCloser(sr io.Reader, sw io.Writer, c io.Closer) secureReadWriteCloser {
	srwc := secureReadWriteCloser{
		sr,
		sw,
		c,
	}
	return srwc
}

// decrypt decrypts bytes using public-key cryptography.
func decrypt(encMsg []byte, sr *secureReader) ([]byte, error) {

	// Read the nonce.
	var nonce [24]byte
	n, err := io.ReadFull(sr.r, nonce[:])
	if err != nil {
		return nil, err
	}

	// Overhead is the number of bytes of overhead when
	// encrypting a message therefor we need to create a
	// buffer larger than just the encrypted message.
	msg := make([]byte, len(encMsg)+box.Overhead)

	// Read the encrypted message.
	n, err = sr.r.Read(msg)
	if err != nil {
		return nil, err
	}

	// Decrypt the message.
	dec, ok := box.Open(encMsg[:0], msg[:n], &nonce, sr.pub, sr.priv)
	if !ok {
		return nil, err
	}

	return dec, nil
}

// encrypt encrypts bytes using public-key cryptography.
func encrypt(msg []byte, rPubKey, sPrivKey *[32]byte) ([]byte, error) {

	// create a nonce to encrypt with.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	// encrypt the message which appends the result to the nonce
	// in order to store the nonce with the encrypted message so that we
	// can use the same nonce when its decrypted.
	encrypted := box.Seal(nonce[:], msg, &nonce, rPubKey, sPrivKey)
	return encrypted, nil
}

// reveiveKey is the receiveing side of a public key exchange over a
// connection.
func receiveKey(conn net.Conn) (*[32]byte, error) {
	pubKey := make([]byte, 32)
	if _, err := conn.Read(pubKey); err != nil {
		return nil, err
	}

	// After reading pubKey bytes, format the pubKey
	// into the desired array of 32 bytes format.
	var fmtKey [32]byte
	copy(fmtKey[:], pubKey)
	return &fmtKey, nil
}

// Perform a pubic key exchange over the connection.
func handshake(conn net.Conn, sKey *[32]byte) (*[32]byte, error) {

	// Send the public key.
	if _, err := conn.Write(sKey[:]); err != nil {
		return nil, err
	}

	// Recieve the public key from the other side
	// of the connection.
	rKey, err := receiveKey(conn)
	if err != nil {
		return nil, err
	}
	return rKey, nil
}

// Dial connects to the server over a secure connections which encrypts and
// decrypts all read and write bytes.
func Dial(addr string) (io.ReadWriteCloser, error) {

	// Connect to server.
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Generate a public/private encryption keypair for the client.
	pub, clientPrivKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Perform a public encryption key exchange to
	// get the servers' public encryption key.
	serverPubKey, err := handshake(conn, pub)
	if err != nil {
		return nil, err
	}

	// Wrap the connection with a secure reader and secure writer
	// which will encrypt/decrypt all read/writes.
	secureConn := newSecureReadWriteCloser(
		NewSecureReader(conn, clientPrivKey, serverPubKey),
		NewSecureWriter(conn, clientPrivKey, serverPubKey),
		conn,
	)

	return secureConn, nil
}

// Serve starts a secure echo server on the given listener.
func Serve(l net.Listener) error {

	// Accept a connection from a client.
	conn, err := l.Accept()
	if err != nil {
		return fmt.Errorf("Error l.Accept : %v", err)
	}
	defer conn.Close()

	// Generate a pub/private keypair.
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("Error Serve GenKey: %v", err)
	}

	// Perform a public encryption key exchange to
	// get the client's public encryption key.
	clientPubKey, err := handshake(conn, pub)
	if err != nil {
		return fmt.Errorf("Error Serve handshake: %v", err)
	}

	sw := NewSecureWriter(conn, priv, clientPubKey)
	sr := NewSecureReader(conn, priv, clientPubKey)

	// Echo the data back to the client over a secure
	// connection that encrypts and decrypts all reads/writes.
	_, err = io.Copy(sw, sr)
	if err != nil {
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
