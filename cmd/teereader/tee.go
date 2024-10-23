package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type progress struct {
	total uint64
}

// Write updates total bytes downloaded and display progress
func (p *progress) Write(b []byte) (int, error) {
	p.total += uint64(len(b))
	fmt.Printf("Downloaded %d bytes...\n", p.total)
	return len(b), nil
}

func main() {
	// URL of the file to be downloaded
	url := "http://storage.googleapis.com/books/ngrams/books/googlebooks-eng-all-5gram-20120701-0.gz"

	// Create HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Create local file to save the downloaded content
	localFile, err := os.Create("downloaded_file.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer localFile.Close()

	// Create a hasher for calculating checksum
	hasher := sha256.New()

	// Create a TeeReader: reads from response body and writes to progress and hasher
	teeReader := io.TeeReader(res.Body, io.MultiWriter(hasher, &progress{}))

	// Write the content from teeReader to localFile
	if _, err := io.Copy(localFile, teeReader); err != nil {
		log.Fatal(err)
	}

	// Print the calculated checksum
	fmt.Printf("File checksum: %x\n", hasher.Sum(nil))
}
