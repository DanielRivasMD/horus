////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func TestPropagateErr_Nil(t *testing.T) {
	if got := PropagateErr("op", "cat", "msg", nil, map[string]any{"k": "v"}); got != nil {
		t.Errorf("PropagateErr with nil err = %v; want nil", got)
	}
}

func TestPropagateErr_PlainError(t *testing.T) {
	plain := errors.New("plain")
	details := map[string]any{"d1": 1, "d2": "two"}

	e := PropagateErr("Op1", "Cat1", "Message1", plain, details)
	h, ok := AsHerror(e)
	if !ok {
		t.Fatalf("PropagateErr did not return an Herror, got %T", e)
	}

	if h.Op != "Op1" {
		t.Errorf("Op = %q; want %q", h.Op, "Op1")
	}
	if h.Category != "Cat1" {
		t.Errorf("Category = %q; want %q", h.Category, "Cat1")
	}
	if h.Message != "Message1" {
		t.Errorf("Message = %q; want %q", h.Message, "Message1")
	}
	if h.Err != plain {
		t.Errorf("Err = %v; want original plain error", h.Err)
	}
	if len(h.Details) != len(details) {
		t.Errorf("Details len = %d; want %d", len(h.Details), len(details))
	}
	for k, v := range details {
		if h.Details[k] != v {
			t.Errorf("Details[%q] = %v; want %v", k, h.Details[k], v)
		}
	}
	if len(h.Stack) == 0 {
		t.Error("Stack trace should be captured")
	}
}

func TestPropagateErr_PreserveBaseCategoryAndMergeDetails(t *testing.T) {
	// Create an existing Herror
	base := NewCategorizedHerror(
		"origOp",
		"baseCat",
		"origMsg",
		errors.New("inner"),
		map[string]any{"a": 1},
	)
	// Propagate without category override, add new detail b=2
	e := PropagateErr("newOp", "", "newMsg", base, map[string]any{"b": 2})
	h, ok := AsHerror(e)
	if !ok {
		t.Fatalf("expected Herror, got %T", e)
	}

	if h.Op != "newOp" {
		t.Errorf("Op = %q; want %q", h.Op, "newOp")
	}
	// category should come from baseCat
	if h.Category != "baseCat" {
		t.Errorf("Category = %q; want %q", h.Category, "baseCat")
	}
	if h.Message != "newMsg" {
		t.Errorf("Message = %q; want %q", h.Message, "newMsg")
	}
	// Underlying error should be the base Herror
	if h.Err != base {
		t.Errorf("Err = %v; want base Herror", h.Err)
	}
	// Details should include both a=1 and b=2
	if len(h.Details) != 2 {
		t.Errorf("Details len = %d; want 2", len(h.Details))
	}
	if h.Details["a"] != 1 {
		t.Errorf("Details[%q] = %v; want %v", "a", h.Details["a"], 1)
	}
	if h.Details["b"] != 2 {
		t.Errorf("Details[%q] = %v; want %v", "b", h.Details["b"], 2)
	}
}

func TestPropagateErr_OverrideCategoryAndDetailCollision(t *testing.T) {
	// Base Herror with category "oldCat" and detail x=1
	base := NewCategorizedHerror("o", "oldCat", "", nil, map[string]any{"x": 1, "y": "z"})

	// Propagate with category override "newCat" and new detail x=2
	e := PropagateErr("op2", "newCat", "msg2", base, map[string]any{"x": 2})
	h, ok := AsHerror(e)
	if !ok {
		t.Fatalf("expected Herror, got %T", e)
	}

	if h.Category != "newCat" {
		t.Errorf("Category = %q; want %q", h.Category, "newCat")
	}
	// x should be overridden
	if h.Details["x"] != 2 {
		t.Errorf("Details[%q] = %v; want %v", "x", h.Details["x"], 2)
	}
	// y should be preserved
	if h.Details["y"] != "z" {
		t.Errorf("Details[%q] = %v; want %v", "y", h.Details["y"], "z")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
