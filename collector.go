package horus

import (
	"bytes"
)

// CollectingError implements both the io.Writer and error interfaces.
// It accumulates written output, which you can later retrieve as an error message.
type CollectingError struct {
	buf bytes.Buffer
}

// Write appends p to the buffer.
func (ce *CollectingError) Write(p []byte) (n int, err error) {
	return ce.buf.Write(p)
}

// Error returns the accumulated message as a string.
func (ce *CollectingError) Error() string {
	return ce.buf.String()
}

// NewCollectingError creates and returns a new CollectingError.
func NewCollectingError() *CollectingError {
	return &CollectingError{}
}
