////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"io"
	"os"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// exitFunc is called by CheckErr to terminate the process.
// You can override this in tests to capture the exit code.
var exitFunc = os.Exit

////////////////////////////////////////////////////////////////////////////////////////////////////

type CheckConfig struct {
	Out       io.Writer
	ExitCode  int
	Formatter FormatterFunc
}

var exit = os.Exit

////////////////////////////////////////////////////////////////////////////////////////////////////

// Global error registry to track error types (by category).
var errorTypeRegistry = make(map[string]int)

////////////////////////////////////////////////////////////////////////////////////////////////////

// RegisterError increments the count of errors for a given category.
func RegisterError(err error) {
	if err == nil {
		return
	}
	if herr, ok := AsHerror(err); ok && herr.Category != "" {
		errorTypeRegistry[herr.Category]++
	} else {
		errorTypeRegistry["unknown"]++
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// GetErrorRegistry returns a copy of the error type registry.
func GetErrorRegistry() map[string]int {
	copyRegistry := make(map[string]int)
	for key, count := range errorTypeRegistry {
		copyRegistry[key] = count
	}
	return copyRegistry
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// checkOpt is a functional option that mutates checkParams.
type checkOpt func(*checkParams)

// checkParams holds all customizable fields for CheckErr.

type checkParams struct {
	op       string
	category string
	message  string
	details  map[string]any
	writer   io.Writer
	exitCode int
}

// WithOp overrides the default operation name.
func WithOp(opName string) checkOpt {
	return func(p *checkParams) {
		p.op = opName
	}
}

// WithCategory overrides the default error category.
func WithCategory(cat string) checkOpt {
	return func(p *checkParams) {
		p.category = cat
	}
}

// WithMessage overrides the default user-facing message.
func WithMessage(msg string) checkOpt {
	return func(p *checkParams) {
		p.message = msg
	}
}

// WithDetails overrides or augments the default metadata map.
func WithDetails(d map[string]any) checkOpt {
	return func(p *checkParams) {
		p.details = d
	}
}

// WithWriter lets you redirect the output (for tests, etc).
func WithWriter(w io.Writer) checkOpt {
	return func(c *checkParams) {
		c.writer = w
	}
}

// WithExitCode lets you override the exit code.
func WithExitCode(code int) checkOpt {
	return func(c *checkParams) {
		c.exitCode = code
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckErr registers, wraps, formats and logs a fatal error via Horus.
// On non-nil err it exits the process with code (default: 1).
func CheckErr(err error, opts ...checkOpt) {
	if err == nil {
		return
	}

	// 1) Metrics, instrumentation
	RegisterError(err)

	// 2) Default parameters
	cfg := checkParams{
		op:       "check error",
		category: "runtime_error",
		message:  "An error occurred during execution",
		details: map[string]any{
			"severity": "critical",
			"location": "checkErr",
		},
		writer:   os.Stderr,
		exitCode: 1,
	}

	// 3) Apply any overrides
	for _, opt := range opts {
		opt(&cfg)
	}

	// 4) Build the rich Herror
	herr := NewCategorizedHerror(
		cfg.op,
		cfg.category,
		cfg.message,
		err,
		cfg.details,
	)

	// 5) Format & print
	//    We know NewCategorizedHerror always returns *Herror under the hood,
	//    but we still AsHerror for safety.
	if he, ok := AsHerror(herr); ok {
		fmt.Fprintln(cfg.writer, PseudoJSONFormatter(he))
	} else {
		// fallback, should never happen
		fmt.Fprintln(cfg.writer, herr.Error())
	}

	// 6) Exit
	exitFunc(cfg.exitCode)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
