////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestIsHerrorAndAsHerror(t *testing.T) {
	plain := errors.New("plain")
	if IsHerror(plain) {
		t.Error("IsHerror(plain) should be false")
	}
	if h, ok := AsHerror(plain); ok || h != nil {
		t.Errorf("AsHerror(plain) = %v, %v; want nil,false", h, ok)
	}

	// wrap a plain error into an Herror
	wrapped := NewHerror("opX", "msgX", plain, nil)
	if !IsHerror(wrapped) {
		t.Error("IsHerror(Herror) should be true")
	}
	h2, ok2 := AsHerror(wrapped)
	if !ok2 {
		t.Error("AsHerror(Herror) ok should be true")
	}
	if h2.Op != "opX" || h2.Message != "msgX" {
		t.Errorf("AsHerror returned wrong Herror: %+v", h2)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestOperationAndUserMessage(t *testing.T) {
	plain := errors.New("nope")
	if op := Operation(plain); op != "" {
		t.Errorf("Operation(plain) = %q; want empty", op)
	}
	if um := UserMessage(plain); um != "" {
		t.Errorf("UserMessage(plain) = %q; want empty", um)
	}

	he := NewHerror("myOp", "myMsg", nil, nil)
	if op := Operation(he); op != "myOp" {
		t.Errorf("Operation(Herror) = %q; want %q", op, "myOp")
	}
	if um := UserMessage(he); um != "myMsg" {
		t.Errorf("UserMessage(Herror) = %q; want %q", um, "myMsg")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestDetailAndAllDetails(t *testing.T) {
	plain := errors.New("oops")
	if v, ok := Detail(plain, "k"); ok || v != nil {
		t.Errorf("Detail(plain) = %v, %v; want nil,false", v, ok)
	}
	if all := AllDetails(plain); all != nil {
		t.Errorf("AllDetails(plain) = %v; want nil", all)
	}

	base := NewHerror("op", "", nil, map[string]any{"x": 42})
	if v, ok := Detail(base, "x"); !ok || v.(int) != 42 {
		t.Errorf("Detail(Herror, \"x\") = %v, %v; want 42,true", v, ok)
	}
	if _, ok := Detail(base, "missing"); ok {
		t.Error("Detail should be false for missing key")
	}
	all := AllDetails(base)
	if all == nil || all["x"].(int) != 42 {
		t.Errorf("AllDetails = %v; want map[x:42]", all)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestCategory(t *testing.T) {
	plain := errors.New("err")
	if cat := Category(plain); cat != "" {
		t.Errorf("Category(plain) = %q; want empty", cat)
	}
	he := NewCategorizedHerror("op", "cat42", "", nil, nil)
	if cat := Category(he); cat != "cat42" {
		t.Errorf("Category(Herror) = %q; want %q", cat, "cat42")
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
