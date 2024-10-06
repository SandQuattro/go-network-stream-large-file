package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"go-network-stream-large-file/proto"
	"go-network-stream-large-file/sender"
	"io"
	"log"
	"net"
	"time"
)

type FileServer struct {
	proto int
}

func main() {
	p := proto.UDP
	go func() {
		time.Sleep(1 * time.Second)
		err := sender.SendFile(p, 1<<20)
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	server := &FileServer{
		proto: p,
	}

	server.start()
}

func (fs *FileServer) start() {
	switch {
	case fs.proto == proto.TCP:
		fs.TCPReader()
	case fs.proto == proto.UDP:
		fs.UDPReader()
	default:
		log.Fatalf("Unsupported protocol: %d", fs.proto)
	}
}

func (fs *FileServer) TCPReader() {
	conn, e := net.Listen(proto.String(fs.proto), ":3000")
	if e != nil {
		log.Fatal(e)
	}
	defer conn.Close()

	for {
		socket, err := conn.Accept()
		if err != nil {
			log.Panic(err)
		}
		go fs.tcpReadLoop(socket)
	}
}

func (fs *FileServer) UDPReader() {
	lAddr := &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 3000,
	}

	conn, err := net.ListenUDP(proto.String(fs.proto), lAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	var totalSize, currentPos int64
	// first we are reading our data actual size from connection
	err = binary.Read(conn, binary.LittleEndian, &totalSize)
	if err != nil {
		log.Panic(err)
	}

	// we are adding checksum of our data to check integrity on receiver
	var checkSum [32]byte
	err = binary.Read(conn, binary.LittleEndian, &checkSum)
	if err != nil {
		log.Panic(err)
	}

	buf := new(bytes.Buffer)
	b := make([]byte, proto.MaxPacketSize)
	n := 0

	for {
		n, err = conn.Read(b)
		if err != nil {
			log.Panic(err)
		}

		currentPos += int64(n)
		buf.Write(b[:n])
		log.Printf("received %d bytes over the network, current pos:%d", n, currentPos)

		if currentPos >= totalSize {
			break
		}
	}

	sum256 := sha256.Sum256(buf.Bytes()[:totalSize])
	if !bytes.Equal(sum256[:], checkSum[:]) {
		log.Panic("!!! checksum mismatch")
	}

	log.Println("checksum correct")

}

func (fs *FileServer) tcpReadLoop(conn io.Reader) {
	buf := new(bytes.Buffer)
	n := int64(0)

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
		n, err = io.CopyN(buf, conn, size)
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
