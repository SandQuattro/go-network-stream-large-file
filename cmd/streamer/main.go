package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"go-network-stream-large-file/sender"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct{}

func (fs *FileServer) start() {
	socket, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Panic(err)
		}
		go fs.readLoop(conn)
	}
}

func main() {
	go func() {
		// schedule us please
		time.Sleep(1 * time.Second)
		err := sender.SendFile(1 << 20)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	server := &FileServer{}
	server.start()
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)

	var size int64
	// first we are reading our data actual size from connection
	err := binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		log.Fatal(err)
	}

	// we are adding checksum of our data to check integrity on receiver
	var checkSum [32]byte
	err = binary.Read(conn, binary.LittleEndian, &checkSum)
	if err != nil {
		log.Fatal(err)
	}

	for {
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}

		// log.Println(buf.Bytes())
		log.Printf("received %d bytes over the network", n)

		sum256 := sha256.Sum256(buf.Bytes())
		if !bytes.Equal(sum256[:], checkSum[:]) {
			log.Fatal("checksum mismatch")
		}

		log.Println("checksum correct")
	}
}
