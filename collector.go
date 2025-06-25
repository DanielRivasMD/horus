////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////

// CollectingError implements both io.Writer and error, accumulating writes
// into an internal buffer. It is safe for concurrent use.
type CollectingError struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

// NewCollectingError returns an empty CollectingError.
func NewCollectingError() *CollectingError {
	return &CollectingError{}
}

// Write appends p to the internal buffer. It is safe for concurrent calls.
func (ce *CollectingError) Write(p []byte) (n int, err error) {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	return ce.buf.Write(p)
}

// WriteString appends s to the internal buffer. It is safe for concurrent calls.
func (ce *CollectingError) WriteString(s string) (n int, err error) {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	return ce.buf.WriteString(s)
}

// Error returns the accumulated contents as a string. It is safe for
// concurrent calls.
func (ce *CollectingError) Error() string {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	return ce.buf.String()
}

// Bytes returns a copy of the internal buffer's bytes. It is safe for
// concurrent calls.
func (ce *CollectingError) Bytes() []byte {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	// Return a copy to avoid exposing underlying buffer
	b := make([]byte, ce.buf.Len())
	copy(b, ce.buf.Bytes())
	return b
}

// Reset clears the internal buffer, allowing reuse of the CollectingError.
// It is safe for concurrent calls.
func (ce *CollectingError) Reset() {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	ce.buf.Reset()
}

////////////////////////////////////////////////////////////////////////////////////////////////////
