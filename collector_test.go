////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
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

////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestCollectingError_MultipleWrites(t *testing.T) {
	ce := NewCollectingError()
	ce.Write([]byte("foo"))
	ce.Write([]byte("bar"))
	if got := ce.Error(); got != "foobar" {
		t.Errorf("multiple writes, Error() = %q; want %q", got, "foobar")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestCollectingError_CompatibleWithFprintf(t *testing.T) {
	ce := NewCollectingError()
	// fmt.Fprintf writes through io.Writer
	fmt.Fprintf(ce, "value=%d; ok=%t\n", 42, true)
	want := "value=42; ok=true\n"
	if got := ce.Error(); got != want {
		t.Errorf("after Fprintf, Error() = %q; want %q", got, want)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
