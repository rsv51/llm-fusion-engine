package util

import (
	"bytes"
	"io"
)

// TeeReader is similar to io.TeeReader, but it captures the entire stream
// into a buffer while it's being read.
type TeeReader struct {
	reader io.Reader
	buffer *bytes.Buffer
}

// NewTeeReader creates a new TeeReader.
func NewTeeReader(reader io.Reader) *TeeReader {
	return &TeeReader{
		reader: reader,
		buffer: new(bytes.Buffer),
	}
}

// Read reads from the original reader and writes to the internal buffer.
func (t *TeeReader) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)
	if n > 0 {
		t.buffer.Write(p[:n])
	}
	return
}

// GetContent returns the full content that has been read so far.
func (t *TeeReader) GetContent() []byte {
	return t.buffer.Bytes()
}

// Close closes the underlying reader if it's an io.ReadCloser.
func (t *TeeReader) Close() error {
	if closer, ok := t.reader.(io.ReadCloser); ok {
		return closer.Close()
	}
	return nil
}