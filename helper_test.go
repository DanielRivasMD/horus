////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestAsHerrorAndOperation(t *testing.T) {
	plain := errors.New("plain")

	// plain errors → AsHerror=false, Operation=false
	if h, ok := AsHerror(plain); ok || h != nil {
		t.Errorf("AsHerror(plain) = %v, %v; want nil,false", h, ok)
	}
	if op, ok := Operation(plain); ok || op != "" {
		t.Errorf("Operation(plain) = %q, %v; want \"\", false", op, ok)
	}

	// wrap a plain error into an Herror
	wrapped := NewHerror("opX", "msgX", plain, nil)

	// AsHerror on Herror → should succeed
	h2, ok2 := AsHerror(wrapped)
	if !ok2 {
		t.Error("AsHerror(Herror) ok should be true")
	}
	if h2.Op != "opX" || h2.Message != "msgX" {
		t.Errorf("AsHerror returned wrong Herror: %+v", h2)
	}

	// Operation on Herror → should return the op and true
	if op, ok := Operation(wrapped); !ok || op != "opX" {
		t.Errorf("Operation(Herror) = %q, %v; want %q, true", op, ok, "opX")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestOperationAndUserMessage(t *testing.T) {
	plain := errors.New("nope")

	// plain errors should return "", false
	if op, ok := Operation(plain); ok || op != "" {
		t.Errorf("Operation(plain) = %q, %v; want \"\", false", op, ok)
	}
	if um := UserMessage(plain); um != "" {
		t.Errorf("UserMessage(plain) = %q; want empty", um)
	}

	// Herror should return its Op and true
	he := NewHerror("myOp", "myMsg", nil, nil)
	op, ok := Operation(he)
	if !ok || op != "myOp" {
		t.Errorf("Operation(Herror) = %q, %v; want %q, true", op, ok, "myOp")
	}
	if um := UserMessage(he); um != "myMsg" {
		t.Errorf("UserMessage(Herror) = %q; want %q", um, "myMsg")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestGetDetailAndDetails(t *testing.T) {
	plain := errors.New("oops")

	// GetDetail on a plain error
	if v, ok := GetDetail(plain, "k"); ok || v != nil {
		t.Errorf("GetDetail(plain) = %v, %v; want nil,false", v, ok)
	}

	// Details should return an empty (but non-nil) map
	all := Details(plain)
	if all == nil {
		t.Errorf("Details(plain) = %v; want non-nil empty map", all)
	}
	if len(all) != 0 {
		t.Errorf("Details(plain) = %v; want empty map", all)
	}

	// Now when Herror has details:
	base := NewHerror("op", "", nil, map[string]any{"x": 42})
	if v, ok := GetDetail(base, "x"); !ok || v.(int) != 42 {
		t.Errorf("GetDetail(Herror, \"x\") = %v, %v; want 42,true", v, ok)
	}
	if _, ok := GetDetail(base, "missing"); ok {
		t.Error("GetDetail should be false for missing key")
	}
	all2 := Details(base)
	if all2 == nil || all2["x"].(int) != 42 {
		t.Errorf("Details(Herror) = %v; want map[x:42]", all2)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestCategory(t *testing.T) {
	plain := errors.New("err")
	// plain errors → empty, false
	if cat, ok := Category(plain); ok || cat != "" {
		t.Errorf("Category(plain) = %q, %v; want \"\", false", cat, ok)
	}

	he := NewCategorizedHerror("op", "cat42", "", nil, nil)
	// Herror → actual category, true
	if cat, ok := Category(he); !ok || cat != "cat42" {
		t.Errorf("Category(Herror) = %q, %v; want %q, true", cat, ok, "cat42")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestStackTraceHelper(t *testing.T) {
	plain := errors.New("err")
	if st := StackTrace(plain); st != "" {
		t.Errorf("StackTrace(plain) = %q; want empty", st)
	}

	// create an Herror so stack is captured
	he := NewHerror("op", "", nil, nil)
	trace := StackTrace(he)
	if trace == "" {
		t.Error("StackTrace(Herror) should not be empty")
	}
	// must mention this test function name somewhere
	if !strings.Contains(trace, "TestStackTraceHelper") {
		t.Errorf("StackTrace missing test frame, got:\n%s", trace)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
