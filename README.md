# Sending and receiving binary data over network with checksum validation

Process:

Generating data:

A file of size bytes is generated first, filled with random data, using a cryptographically secure random number generator (crypto/rand).
Establishing a connection:

A TCP connection is established to the server on port 3000. If the connection fails, the function returns with an error. The connection is automatically closed after the function returns thanks to defer conn.Close().
Transmitting the file size:

The file size (size) is written to the connection using LittleEndian byte order. This lets the receiver know how much data to expect.
Transmitting the SHA-256 hash:

To verify the integrity of the data on the receiving side, a SHA-256 hash is calculated from the file contents. This hash is also written to the connection.
Transmitting data:

The file data is transferred to the connection using the io.CopyN method. This method ensures that all data is transferred until the specified size is reached.
Logging:

Information about the number of bytes transferred is logged.

# Future Plans

- [ ] add udp data transfer
- [x] add transferred data checksum validation