////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// stripANSI removes ANSI color escapes so we can assert on plain text.
var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

////////////////////////////////////////////////////////////////////////////////////////////////////

func stripANSI(s string) string {
	return ansi.ReplaceAllString(s, "")
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestJSONFormatter(t *testing.T) {
	original := &Herror{
		Op:       "myOp",
		Message:  "myMsg",
		Err:      errors.New("oh no"),
		Details:  map[string]any{"foo": "bar"},
		Category: "catX",
		Stack:    []uintptr{1, 2, 3},
	}

	out := JSONFormatter(original)

	// 1) Unmarshal into a map so Err (an interface) doesn't blow up.
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON from JSONFormatter: %v\n%s", err, out)
	}

	// 2) Helper for pulling a string field out of the map
	getString := func(key string) string {
		v, ok := parsed[key].(string)
		if !ok {
			t.Fatalf("field %q is not a string or missing (%T)", key, parsed[key])
		}
		return v
	}

	// Compare core scalar fields
	if got, want := getString("Op"), original.Op; got != want {
		t.Errorf("Op = %q; want %q", got, want)
	}
	if got, want := getString("Message"), original.Message; got != want {
		t.Errorf("Message = %q; want %q", got, want)
	}
	if got, want := getString("Err"), original.Err.Error(); got != want {
		t.Errorf("Err = %q; want %q", got, want)
	}
	if got, want := getString("Category"), original.Category; got != want {
		t.Errorf("Category = %q; want %q", got, want)
	}

	// Compare Details map
	rawDetails, ok := parsed["Details"].(map[string]any)
	if !ok {
		t.Fatalf("Details is not a map[string]any: %T", parsed["Details"])
	}
	if got, want := rawDetails["foo"], original.Details["foo"]; got != want {
		t.Errorf("Details[\"foo\"] = %v; want %v", got, want)
	}

	// Compare Stack length
	rawStack, ok := parsed["Stack"].([]any)
	if !ok {
		t.Fatalf("Stack is not a []any: %T", parsed["Stack"])
	}
	if len(rawStack) != len(original.Stack) {
		t.Errorf("Stack length = %d; want %d", len(rawStack), len(original.Stack))
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestSimpleColoredFormatter(t *testing.T) {
	// 1) Create an Herror and unwrap it
	err := NewCategorizedHerror("op", "cat", "msg", errors.New("inner"), nil)
	h, ok := AsHerror(err)
	if !ok {
		t.Fatalf("expected an Herror, got %T", err)
	}

	// 2) Generate colored output, then strip ANSI codes
	raw := SimpleColoredFormatter(h)
	plain := stripANSI(raw)

	// 3) Verify the plain‐text prefix
	wantPrefix := "ERROR: " + h.Error()
	if !strings.HasPrefix(plain, wantPrefix) {
		t.Errorf("SimpleColoredFormatter = %q; want prefix %q", plain, wantPrefix)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestPseudoJSONFormatter_Content(t *testing.T) {
	// Build a minimal Herror
	h := &Herror{
		Op:       "fooOp",
		Message:  "fooMsg",
		Err:      errors.New("error123"),
		Details:  map[string]any{"key1": "val1"},
		Category: "myCat",
		Stack:    []uintptr{}, // ignore the actual frames here
	}

	out := stripANSI(PseudoJSONFormatter(h))
	lines := strings.Split(out, "\n")

	// Each row says “key” must appear, “value” (if any) must appear,
	// and the key must come before the value on that line.
	tests := []struct {
		key   string
		value string
	}{
		{"Op", "fooOp,"},
		{"Message", "fooMsg,"},
		{"Err", "error123,"},
		{"Details", ""},   // just the header
		{"key1", "val1,"}, // nested detail
		{"Category", "myCat,"},
		{"Stack", ""}, // just the header
	}

	for _, tc := range tests {
		found := false
		for _, line := range lines {
			if !strings.Contains(line, tc.key) {
				continue
			}
			if tc.value != "" && !strings.Contains(line, tc.value) {
				continue
			}
			// if both key & value, ensure key precedes value
			if tc.value != "" && strings.Index(line, tc.key) > strings.Index(line, tc.value) {
				continue
			}
			found = true
			break
		}
		if !found {
			t.Errorf(
				"PseudoJSONFormatter missing a line with key %q value %q\nFull output:\n%s",
				tc.key, tc.value, out,
			)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestFormatPanic(t *testing.T) {
	raw := FormatPanic("DoIt", "it broke")
	plain := stripANSI(raw) // strip out the \x1b[…m codes

	want := "Panic [DoIt]: it broke"
	if plain != want {
		t.Errorf("FormatPanic = %q; want %q", plain, want)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
