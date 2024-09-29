package main

import (
	"bytes"
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
		time.Sleep(4 * time.Second)
		err := sender.SendFile(100)
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

	for {
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(buf.Bytes())
		log.Printf("received %d bytes over the network", n)
	}
}
