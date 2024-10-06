package sender

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"go-network-stream-large-file/proto"
	"io"
	"log"
	"net"
	"time"
)

func SendFile(p int, size int64) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial(proto.String(p), ":3000")
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

	var n int64
	if p == proto.TCP {
		// actual streaming our data to connection
		n, err = io.CopyN(conn, bytes.NewReader(file), size)
		if err != nil {
			return err
		}
	} else if p == proto.UDP {
		buf := new(bytes.Buffer)

		buf.Write(file)

		for buf.Len() > 0 {
			chunkSize := proto.MaxPacketSize
			if buf.Len() < proto.MaxPacketSize {
				chunkSize = buf.Len()
			}

			chunk := buf.Next(chunkSize)
			num, err := writeData(conn, chunk)
			if err != nil {
				log.Panic(err)
			}
			n += int64(num)

			// we can't do this soo fast!
			time.Sleep(1 * time.Nanosecond)
		}
	}

	// log.Print(file)
	log.Printf("sent %d bytes over connection", n)

	return nil
}

func writeData(dst io.Writer, data []byte) (int, error) {
	return dst.Write(data)
}
