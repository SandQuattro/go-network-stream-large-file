package sender

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"log"
	"net"
)

func SendFile(size int64) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// first we write our data size into connection, and then actual data
	err = binary.Write(conn, binary.LittleEndian, size)
	if err != nil {
		return err
	}

	// we are adding checksum of our data to check integrity on receiver
	sum256 := sha256.Sum256(file)
	err = binary.Write(conn, binary.LittleEndian, sum256)
	if err != nil {
		return err
	}

	// actual streaming our data to connection
	n, err := io.CopyN(conn, bytes.NewReader(file), size)
	if err != nil {
		return err
	}

	// log.Print(file)
	log.Printf("sent %d bytes over connection", n)

	return nil
}
