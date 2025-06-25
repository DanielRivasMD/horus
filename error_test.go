////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// captureOutput temporarily redirects stdout and returns what was printed.
func captureOutput_error(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestHerror_FormatAndUnwrap(t *testing.T) {
	h := &Herror{Op: "X", Message: "M", Err: errors.New("inner")}
	// custom formatter just echoes Op
	f := FormatterFunc(func(h *Herror) string { return "->" + h.Op })
	if got := h.Format(f); got != "->X" {
		t.Errorf("Format() = %q; want %q", got, "->X")
	}
	if got := h.Unwrap(); got.Error() != "inner" {
		t.Errorf("Unwrap() = %v; want inner", got)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestHerror_StackTraceAndHasStack(t *testing.T) {
	// ensure captureStack actually populates
	h := NewHerror("op", "", nil, nil).(*Herror)
	if !h.HasStack() {
		t.Fatal("HasStack should be true after NewHerror")
	}
	trace := h.StackTrace()
	if trace == "" {
		t.Error("StackTrace() returned empty")
	}
	// stack trace must mention this test function
	if !strings.Contains(trace, "TestHerror_StackTraceAndHasStack") {
		t.Errorf("StackTrace() = %q; want to contain test name", trace)
	}

	// empty-out stack
	h2 := &Herror{Op: "o"}
	if h2.HasStack() {
		t.Error("HasStack on zero‐value should be false")
	}
	if st := h2.StackTrace(); st != "" {
		t.Errorf("StackTrace on zero‐stack = %q; want empty", st)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

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
	// Err should be the string
	if m["Err"] != "E" {
		t.Errorf("Err = %v; want %q", m["Err"], "E")
	}
	// check some other fields
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

////////////////////////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestWrap(t *testing.T) {
	// nil in → nil out
	if w := Wrap(nil, "X", "Y"); w != nil {
		t.Error("Wrap(nil) should be nil")
	}
	// plain error
	pe := errors.New("plain")
	w1 := Wrap(pe, "OpW", "MsgW")
	he1, ok := AsHerror(w1)
	if !ok || he1.Err.Error() != "plain" {
		t.Error("Wrap on plain failed:", w1)
	}
	if he1.Op != "OpW" || he1.Message != "MsgW" {
		t.Error("Wrap fields incorrect:", he1)
	}
	// wrapping existing Herror preserves its Details & Category
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

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestWithDetail(t *testing.T) {
	// plain error → new Herror
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

	// existing Herror → augments details in place
	h2 := NewHerror("O", "M", nil, map[string]any{"a": 1})
	d2 := WithDetail(h2, "b", 2)
	he2, _ := AsHerror(d2)
	if he2.Details["a"] != 1 || he2.Details["b"] != 2 {
		t.Errorf("WithDetail did not merge details: %v", he2.Details)
	}
	// pointer identity preserved
	if he2 != h2 {
		t.Error("WithDetail should return same *Herror instance")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestPanicFunc(t *testing.T) {
	// capture printed panic message then recover
	out := captureOutput_error(func() {
		defer func() {
			r := recover()
			he, ok := r.(*Herror)
			if !ok {
				t.Fatalf("expected *Herror, got %T", r)
			}
			if he.Op != "P" || he.Message != "M" {
				t.Errorf("Panic Herror incorrect: %+v", he)
			}
		}()
		Panic("P", "M")
	})
	// FormatPanic wraps in ANSI; strip and assert
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	plain := ansi.ReplaceAllString(out, "")
	if !strings.Contains(plain, "Panic [P]: M") {
		t.Errorf("Panic printed %q; want Panic [P]: M", plain)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
