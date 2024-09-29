package main

import (
	"bytes"
	"crypto/rand"
	"errors"
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

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	buf := new(bytes.Buffer)
	for {
		n, err := io.CopyN(buf, conn, 4000)
		if err != nil {
			log.Fatal(err)
		}
		if errors.Is(err, io.EOF) {
			log.Println("Done reading file")
			break
		}

		log.Println(buf.Bytes())
		log.Printf("received %d bytes over the network", n)
	}
}

func sendFile(size int64) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	n, err := io.CopyN(conn, bytes.NewReader(file), size)
	if err != nil {
		return err
	}

	// log.Print(file)
	log.Printf("sent %d bytes over connection", n)

	return nil
}

func main() {
	go func() {
		time.Sleep(4 * time.Second)
		err := sendFile(4000)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	server := &FileServer{}
	server.start()
}
