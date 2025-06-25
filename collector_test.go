////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"fmt"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestNewCollectingError_Empty(t *testing.T) {
	ce := NewCollectingError()
	if ce == nil {
		t.Fatal("NewCollectingError returned nil")
	}
	if got := ce.Error(); got != "" {
		t.Errorf("empty CollectingError.Error() = %q; want empty string", got)
	}
}

func TestCollectingError_ZeroWrite(t *testing.T) {
	ce := NewCollectingError()
	n, err := ce.Write([]byte{})
	if err != nil {
		t.Fatalf("Write(empty) returned error: %v", err)
	}
	if n != 0 {
		t.Errorf("Write(empty) = %d; want 0", n)
	}
	if got := ce.Error(); got != "" {
		t.Errorf("after zero‐write, Error() = %q; want empty string", got)
	}
}

func TestCollectingError_SingleWrite(t *testing.T) {
	ce := NewCollectingError()
	data := []byte("hello")
	n, err := ce.Write(data)
	if err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if n != len(data) {
		t.Errorf("Write returned n = %d; want %d", n, len(data))
	}
	if got := ce.Error(); got != "hello" {
		t.Errorf("Error() = %q; want %q", got, "hello")
	}
}

func TestCollectingError_MultipleWrites(t *testing.T) {
	ce := NewCollectingError()
	ce.Write([]byte("foo"))
	ce.Write([]byte("bar"))
	if got := ce.Error(); got != "foobar" {
		t.Errorf("multiple writes, Error() = %q; want %q", got, "foobar")
	}
}

func TestCollectingError_CompatibleWithFprintf(t *testing.T) {
	ce := NewCollectingError()
	fmt.Fprintf(ce, "value=%d; ok=%t\n", 42, true)
	want := "value=42; ok=true\n"
	if got := ce.Error(); got != want {
		t.Errorf("after Fprintf, Error() = %q; want %q", got, want)
	}
}

func TestCollectingError_WriteStringAndBytes(t *testing.T) {
	ce := NewCollectingError()
	s := "gopher"
	n, err := ce.WriteString(s)
	if err != nil {
		t.Fatalf("WriteString returned error: %v", err)
	}
	if n != len(s) {
		t.Errorf("WriteString returned n = %d; want %d", n, len(s))
	}

	// Error() should match
	if got := ce.Error(); got != s {
		t.Errorf("Error() = %q; want %q", got, s)
	}

	// Bytes() returns a copy
	b := ce.Bytes()
	if !bytes.Equal(b, []byte(s)) {
		t.Errorf("Bytes() = %q; want %q", b, s)
	}
	// Mutate returned slice and ensure underlying buffer doesn't change
	b[0] = 'X'
	if ce.Error()[0] == 'X' {
		t.Error("Bytes() should return a copy, not the internal buffer")
	}
}

func TestCollectingError_Reset(t *testing.T) {
	ce := NewCollectingError()
	ce.WriteString("first")
	ce.Reset()
	if got := ce.Error(); got != "" {
		t.Errorf("after Reset, Error() = %q; want empty", got)
	}

	// Should be reusable
	ce.WriteString("second")
	if got := ce.Error(); got != "second" {
		t.Errorf("after Reset and write, Error() = %q; want %q", got, "second")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
