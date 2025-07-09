////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// captureOutput_error temporarily redirects stderr and returns what was printed.
func captureOutput_error(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestHerror_Error(t *testing.T) {
	base := &Herror{Op: "foo"}
	if got := base.Error(); got != "operation 'foo' failed" {
		t.Errorf("Error() = %q; want %q", got, "operation 'foo' failed")
	}

	withMsg := &Herror{Op: "foo", Message: "bar"}
	if !strings.Contains(withMsg.Error(), ": bar") {
		t.Errorf("Expected message part in %q", withMsg.Error())
	}

	withErr := &Herror{Op: "foo", Err: errors.New("beep")}
	if !strings.Contains(withErr.Error(), "(caused by: beep)") {
		t.Errorf("Expected cause in %q", withErr.Error())
	}

	withDetails := &Herror{
		Op:      "foo",
		Details: map[string]any{"x": 1},
	}
	if !strings.Contains(withDetails.Error(), "(details: map[x:1])") {
		t.Errorf("Expected details in %q", withDetails.Error())
	}

	full := &Herror{
		Op:       "opA",
		Message:  "msgA",
		Err:      errors.New("errA"),
		Details:  map[string]any{"k": "v"},
		Category: "catA",
	}
	out := full.Error()
	wantPieces := []string{
		"operation 'opA' failed: msgA",
		"(caused by: errA)",
		"(details: map[k:v])",
		"[category: catA]",
	}
	for _, piece := range wantPieces {
		if !strings.Contains(out, piece) {
			t.Errorf("Error() missing %q in %q", piece, out)
		}
	}
}

func TestHerror_UnwrapAndFmt(t *testing.T) {
	// Test Unwrap()
	h := &Herror{Op: "X", Message: "M", Err: errors.New("inner"), Stack: nil}
	if u := h.Unwrap(); u == nil || u.Error() != "inner" {
		t.Errorf("Unwrap() = %v; want inner", u)
	}

	// %v invokes Error()
	want := "operation 'X' failed: M (caused by: inner)"
	if got := fmt.Sprintf("%v", h); got != want {
		t.Errorf("%%v = %q; want %q", got, want)
	}

	// %+v should include stack frames (indented lines)
	h2 := NewHerror("Y", "msg", nil, nil)
	s := fmt.Sprintf("%+v", h2)
	if !strings.Contains(s, "operation 'Y' failed: msg") {
		t.Errorf("%%+v missing error message; got %q", s)
	}
	if !strings.Contains(s, "\n\t") {
		t.Errorf("%%+v missing any stack frames; got %q", s)
	}
}

func TestHerror_MarshalJSON(t *testing.T) {
	h := &Herror{
		Op:       "O",
		Message:  "M",
		Err:      errors.New("E"),
		Details:  map[string]any{"a": 1},
		Category: "C",
		Stack:    []uintptr{0, 1},
	}
	raw, err := json.Marshal(h)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["Err"] != "E" {
		t.Errorf("Err = %v; want %q", m["Err"], "E")
	}
	if m["Op"] != "O" || m["Message"] != "M" || m["Category"] != "C" {
		t.Errorf("JSON fields wrong: %v", m)
	}
	if details, ok := m["Details"].(map[string]any); !ok || details["a"] != float64(1) {
		t.Errorf("Details wrong: %v", m["Details"])
	}
	if stack, ok := m["Stack"].([]any); !ok || len(stack) != 2 {
		t.Errorf("Stack wrong: %v", m["Stack"])
	}
}

func TestNewHerror_Constructors(t *testing.T) {
	plain := errors.New("X")
	h1 := NewHerror("Op1", "Msg1", plain, map[string]any{"d": 2})
	he1, ok := AsHerror(h1)
	if !ok {
		t.Fatalf("NewHerror did not return Herror")
	}
	if he1.Op != "Op1" || he1.Message != "Msg1" || he1.Err != plain {
		t.Error("NewHerror fields incorrect:", he1)
	}
	if he1.Details["d"] != 2 {
		t.Error("NewHerror details incorrect:", he1.Details)
	}

	h2 := NewCategorizedHerror("Op2", "Cat2", "Msg2", nil, nil)
	he2, _ := AsHerror(h2)
	if he2.Category != "Cat2" {
		t.Error("Category not set:", he2.Category)
	}

	h3 := NewHerrorErrorf("Op3", "%s-%d", "A", 7)
	he3, _ := AsHerror(h3)
	if he3.Message != "A-7" {
		t.Errorf("NewHerrorErrorf Message = %q; want %q", he3.Message, "A-7")
	}
}

func TestWrap(t *testing.T) {
	if w := Wrap(nil, "X", "Y"); w != nil {
		t.Error("Wrap(nil) should be nil")
	}

	pe := errors.New("plain")
	w1 := Wrap(pe, "OpW", "MsgW")
	he1, ok := AsHerror(w1)
	if !ok || he1.Err.Error() != "plain" {
		t.Error("Wrap on plain failed:", w1)
	}
	if he1.Op != "OpW" || he1.Message != "MsgW" {
		t.Error("Wrap fields incorrect:", he1)
	}

	base := NewCategorizedHerror("OpB", "CatB", "MsgB", pe, map[string]any{"x": 1})
	w2 := Wrap(base, "Op2", "Msg2")
	he2, _ := AsHerror(w2)
	if he2.Category != "CatB" || he2.Details["x"] != 1 {
		t.Error("Wrap lost metadata:", he2)
	}
	if he2.Op != "Op2" || he2.Message != "Msg2" {
		t.Error("Wrap did not override op/msg:", he2)
	}
}

func TestWithDetail(t *testing.T) {
	pe := errors.New("plainX")
	d1 := WithDetail(pe, "k1", "v1")
	he1, ok := AsHerror(d1)
	if !ok {
		t.Fatal("WithDetail did not wrap plain error")
	}
	if he1.Op != "unknown" || he1.Message != "plainX" {
		t.Error("WithDetail fields incorrect:", he1)
	}
	if he1.Details["k1"] != "v1" {
		t.Error("WithDetail did not set detail:", he1.Details)
	}

	h2 := NewHerror("O", "M", nil, map[string]any{"a": 1})
	d2 := WithDetail(h2, "b", 2)
	he2, _ := AsHerror(d2)
	if he2.Details["a"] != 1 || he2.Details["b"] != 2 {
		t.Errorf("WithDetail did not merge details: %v", he2.Details)
	}
	if he2 != h2 {
		t.Error("WithDetail should return same *Herror instance")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
