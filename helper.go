////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

type HErrorer interface {
	Herror() *Herror
}

func (h *Herror) Herror() *Herror { return h }

////////////////////////////////////////////////////////////////////////////////////////////////////

func IsHerror(err error) bool {
	_, ok := AsHerror(err)
	return ok
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// AsHerror tries to extract an Herror from the error chain.
func AsHerror(err error) (*Herror, bool) {
	// first try the standard library walk
	var target *Herror
	if errors.As(err, &target) {
		return target, true
	}
	// then try the interface
	if richer, ok := err.(HErrorer); ok {
		return richer.Herror(), true
	}
	return nil, false
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// Operation returns the operation associated with an Herror, if present.
func Operation(err error) (string, bool) {
	if herr, ok := AsHerror(err); ok {
		return herr.Op, true
	}
	return "", false
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// UserMessage returns the user-friendly message associated with an Herror, if present.
func UserMessage(err error) (string, bool) {
	if herr, ok := AsHerror(err); ok {
		return herr.Message, true
	}
	return "", false
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// GetDetail returns a specific detail associated with an Herror, if present.
func GetDetail(err error, key string) (any, bool) {
	if herr, ok := AsHerror(err); ok && herr.Details != nil {
		value, exists := herr.Details[key]
		return value, exists
	}
	return nil, false
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// Details returns all details associated with an Herror, if present.
// Details returns all details associated with an Herror, or
// an empty (non-nil) map if there is none.
func Details(err error) map[string]any {
	if h, ok := AsHerror(err); ok {
		if h.Details != nil {
			return h.Details
		}
		return map[string]any{} // non-nil empty
	}
	return map[string]any{} // non-nil empty
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// Category returns the category associated with an Herror, if present.
func Category(err error) (string, bool) {
	if herr, ok := AsHerror(err); ok {
		return herr.Category, true
	}
	return "", false
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// StackTrace returns the formatted stack trace from an error if it's an Herror.
func StackTrace(err error) (string, bool) {
	if herr, ok := AsHerror(err); ok {
		return herr.StackTrace(), true
	}
	return "", false
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func RootCauseHelper(err error) error {
	for {
		if un := errors.Unwrap(err); un != nil {
			err = un
			continue
		}
		return err
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
