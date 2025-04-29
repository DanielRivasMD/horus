////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// Herror represents a generalized error with added context.
type Herror struct {
	Op      string // The operation being performed (e.g., "read file", "connect db")
	Message string // A user-friendly message providing more context
	Err     error  // The underlying error, if any
	Details map[string]interface{} // Optional details for more specific context
}

func (e *Herror) Error() string {
	msg := fmt.Sprintf("operation '%s' failed", e.Op)
	if e.Message != "" {
		msg += fmt.Sprintf(": %s", e.Message)
	}
	if e.Err != nil {
		msg += fmt.Sprintf(" (caused by: %v)", e.Err)
	}
	if len(e.Details) > 0 {
		msg += fmt.Sprintf(" (details: %v)", e.Details)
	}
	return msg
}

// Unwrap provides access to the underlying error.
func (e *Herror) Unwrap() error {
	return e.Err
}

// NewHerror creates a new Herror.
func NewHerror(op string, message string, err error, details map[string]interface{}) error {
	return &Herror{
		Op:      op,
		Message: message,
		Err:     err,
		Details: details,
	}
}

// NewHerrorErrorf creates a new Herror with a formatted message.
func NewHerrorErrorf(op string, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return &Herror{
		Op:      op,
		Message: message,
	}
}

// WithDetail adds a key-value detail to an existing Herror. If the error is not an Herror,
// a new Herror wrapping the original will be returned with the detail.
func WithDetail(err error, key string, value interface{}) error {
	if herr, ok := err.(*Herror); ok {
		if herr.Details == nil {
			herr.Details = make(map[string]interface{})
		}
		herr.Details[key] = value
		return herr
	}
	return &Herror{
		Op:      "unknown", // Or try to infer from err.Error() if needed
		Message: err.Error(),
		Err:     err,
		Details: map[string]interface{}{key: value},
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
