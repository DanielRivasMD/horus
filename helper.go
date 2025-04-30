////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// IsHerror checks if the given error or any error in its chain is an Herror.
func IsHerror(err error) bool {
	var target *Herror
	return errors.As(err, &target)
}

// AsHerror tries to extract an Herror from the error chain.
func AsHerror(err error) (*Herror, bool) {
	var target *Herror
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}

// Operation returns the operation associated with an Herror, if present.
func Operation(err error) string {
	if herr, ok := AsHerror(err); ok {
		return herr.Op
	}
	return ""
}

// UserMessage returns the user-friendly message associated with an Herror, if present.
func UserMessage(err error) string {
	if herr, ok := AsHerror(err); ok {
		return herr.Message
	}
	return ""
}

// Detail returns a specific detail associated with an Herror, if present.
func Detail(err error, key string) (any, bool) {
	if herr, ok := AsHerror(err); ok && herr.Details != nil {
		value, exists := herr.Details[key]
		return value, exists
	}
	return nil, false
}

// AllDetails returns all details associated with an Herror, if present.
func AllDetails(err error) map[string]any {
	if herr, ok := AsHerror(err); ok {
		return herr.Details
	}
	return nil
}

// Category returns the category associated with an Herror, if present.
func Category(err error) string {
	if herr, ok := AsHerror(err); ok {
		return herr.Category
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////////////////////////
