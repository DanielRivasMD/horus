////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// Herror represents a generalized error with added context.
type Herror struct {
	Op       string         // The operation being performed (e.g., "read file", "connect db")
	Message  string         // A user-friendly message providing more context
	Err      error          // The underlying error, if any
	Details  map[string]any // Optional details for more specific context
	Category string         // Error category (e.g., validation, IO, etc.)
	Stack    []uintptr      // Stack trace captured at the time of error creation.
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

////////////////////////////////////////////////////////////////////////////////////////////////////

// Format generates a custom representation of the error using a formatter function.
func (e *Herror) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "%s\n%s", e.Error(), e.StackTrace())
			return
		}
	}
	io.WriteString(f, e.Error())
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// Unwrap provides access to the underlying error.
func (e *Herror) Unwrap() error {
	return e.Err
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// StackTrace returns a formatted stack trace captured when the error was created.
func (e *Herror) StackTrace() string {
	if len(e.Stack) == 0 {
		return ""
	}
	frames := runtime.CallersFrames(e.Stack)
	var sb strings.Builder
	for {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return sb.String()
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func (h *Herror) HasStack() bool { return len(h.Stack) > 0 }

////////////////////////////////////////////////////////////////////////////////////////////////////

// MarshalJSON ensures Err is emitted as its Error() string, not an object.
func (h *Herror) MarshalJSON() ([]byte, error) {
	type alias Herror
	// if there’s no inner error, marshal it as empty string
	errMsg := ""
	if h.Err != nil {
		errMsg = h.Err.Error()
	}
	return json.Marshal(&struct {
		Err string `json:"Err"`
		*alias
	}{
		Err:   errMsg,
		alias: (*alias)(h),
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// captureStack captures the current call stack.
func captureStack() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	return pcs[:n]
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func newHerror(
	op, category, message string,
	err error,
	details map[string]any,
) *Herror {
	if details == nil {
		details = make(map[string]any)
	}
	return &Herror{
		Op:       op,
		Message:  message,
		Err:      err,
		Details:  details,
		Category: category,
		Stack:    captureStack(),
	}
}

// NewHerror creates a new Herror (no category).
func NewHerror(op, msg string, err error, details map[string]any) error {
	return newHerror(op, "", msg, err, details)
}

// NewCategorizedHerror creates a new Herror with a category.
func NewCategorizedHerror(
	op, category, msg string,
	err error,
	details map[string]any,
) error {
	return newHerror(op, category, msg, err, details)
}

// NewHerrorErrorf creates a new Herror with a formatted message.
func NewHerrorErrorf(op, fmtStr string, args ...any) error {
	return newHerror(op, "", fmt.Sprintf(fmtStr, args...), nil, nil)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// Wrap wraps an existing error with additional context, capturing the stack trace.
func Wrap(err error, op, message string) error {
	if err == nil {
		return nil
	}
	if herr, ok := AsHerror(err); ok {
		return &Herror{
			Op:       op,
			Message:  message,
			Err:      err,
			Details:  herr.Details,
			Category: herr.Category,
			Stack:    captureStack(),
		}
	}
	return &Herror{
		Op:      op,
		Message: message,
		Err:     err,
		Stack:   captureStack(),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// WithDetail adds a key-value detail to an existing Herror. If the error is not an Herror,
// a new Herror wrapping the original will be returned with the detail.
func WithDetail(err error, k string, v any) error {
	if herr, ok := AsHerror(err); ok {
		copy := make(map[string]any, len(herr.Details)+1)
		for kk, vv := range herr.Details {
			copy[kk] = vv
		}
		copy[k] = v
		herr.Details = copy
		return herr
	}
	return newHerror("unknown", "", err.Error(), err, map[string]any{k: v})
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// Panic creates an Herror with context and a captured stack trace, logs the formatted message,
// and then panics with the generated Herror.
func Panic(op, msg string) {
	herr := newHerror(op, "", msg, nil, nil)
	fmt.Fprintln(os.Stderr, FormatPanic(op, msg))
	panic(herr)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func RootCause(err error) error {
	for {
		next := errors.Unwrap(err)
		if next == nil {
			return err
		}
		err = next
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
