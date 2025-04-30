////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// Herror represents a generalized error with added context.
type Herror struct {
	Op       string         // The operation being performed (e.g., "read file", "connect db")
	Message  string         // A user-friendly message providing more context
	Err      error          // The underlying error, if any
	Details  map[string]any // Optional details for more specific context
	Category string         // Error category (e.g., validation, IO, etc.)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// Error generates a human-readable representation of the error.
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
	if e.Category != "" {
		msg += fmt.Sprintf(" [category: %s]", e.Category)
	}
	return msg
}

// Format generates a custom representation of the error using a formatter function.
func (e *Herror) Format(formatter FormatterFunc) string {
	return formatter(e)
}

// Unwrap provides access to the underlying error.
func (e *Herror) Unwrap() error {
	return e.Err
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// NewHerror creates a new Herror.
func NewHerror(op, message string, err error, details map[string]any) error {
	return &Herror{
		Op:      op,
		Message: message,
		Err:     err,
		Details: details,
	}
}

// NewCategorizedHerror creates a new Herror with an error category.
func NewCategorizedHerror(op, category, message string, err error, details map[string]any) error {
	return &Herror{
		Op:       op,
		Message:  message,
		Err:      err,
		Details:  details,
		Category: category,
	}
}

// NewHerrorErrorf creates a new Herror with a formatted message.
func NewHerrorErrorf(op string, format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	return &Herror{
		Op:      op,
		Message: message,
	}
}

// WithDetail adds a key-value detail to an existing Herror. If the error is not an Herror,
// a new Herror wrapping the original will be returned with the detail.
func WithDetail(err error, key string, value any) error {
	if herr, ok := err.(*Herror); ok {
		if herr.Details == nil {
			herr.Details = make(map[string]any)
		}
		herr.Details[key] = value
		return herr
	}
	return &Herror{
		Op:      "unknown",
		Message: err.Error(),
		Err:     err,
		Details: map[string]any{key: value},
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
