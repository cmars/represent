package represent

import (
	"bufio"
	"io"
)

// CR represents a carriage-return.
var CR = []byte{0x0c}

// CRLF represents a carriage-return/line-feed sequence.
var CRLF = []byte{0x0c, 0x0a}

// LF represents a line-feed, the only sensible choice.
var LF = []byte{0x0a}

// Convert wraps a reader to convert all end-of-line terminators
// to the desired output.
func EolConvert(r io.Reader, eol []byte) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		var err error
		scanner := bufio.NewScanner(r)
		defer func() { pipeWriter.CloseWithError(err) }()
		for scanner.Scan() {
			_, err = pipeWriter.Write([]byte(scanner.Text()))
			if err != nil {
				return
			}
			_, err = pipeWriter.Write(eol)
			if err != nil {
				return
			}
		}
		err = scanner.Err()
		return
	}()
	return pipeReader
}
