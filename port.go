package erlang

import (
	"io"
	"os"
)

// Port is used to implement io.ReadWriter for communication with an Erlang port
// interface. This has some minor differences compared to just communicating
// over stdin and stdout. We use this to abstract over all that.
//
// Port implements io.Writer, but not io.Reader (because we probably would not
// know the length of the message we are reading beforehand).
type Port struct {
	reader io.Reader
	writer io.Writer
}

// NewPort is used to return a new port with the given reader
func NewPort() *Port {
	return &Port{reader: os.Stdin, writer: os.Stdout}
}

// ReadMsg is a function used to test reading the input from the given reader.
// Erlang encodes the length into the first two bytes using a binary format. We
// can use that length to determine how big the body is, and read just that.
func (p Port) ReadMsg() ([]byte, error) {
	sizeHeader := make([]byte, 2)
	_, err := p.reader.Read(sizeHeader)

	if err != nil {
		return nil, err
	}
	size := decodeSize(sizeHeader)
	body := make([]byte, size)
	_, err = p.reader.Read(body)

	return body, err
}

// Write is almost a reverse of ReadMsg. We need to encode the length into
// the first two bytes, and then put the remainder of the body in after that.
// Incidentally, this also implements the io.Writer interface.
func (p Port) Write(msg []byte) (n int, err error) {
	length := len(msg)
	sizeHeader := encodeSize(length)
	body := append(sizeHeader, msg...)

	return p.writer.Write(body)
}

func decodeSize(header []byte) int {
	return int(header[0])<<8 | int(header[1])
}

func encodeSize(length int) []byte {
	return []byte{byte(length >> 8 & 0xff), byte(length & 0xff)}
}
