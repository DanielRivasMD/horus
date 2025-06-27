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
// Override in tests if you want to capture the exit code.
var exitFunc = os.Exit

// errorTypeRegistry tracks how many errors of each Category have been seen.
var errorTypeRegistry = make(map[string]int)

////////////////////////////////////////////////////////////////////////////////////////////////////

// RegisterError increments the count for this error’s category.
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

// GetErrorRegistry returns a copy of the current error counts.
func GetErrorRegistry() map[string]int {
	copyMap := make(map[string]int, len(errorTypeRegistry))
	for k, v := range errorTypeRegistry {
		copyMap[k] = v
	}
	return copyMap
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// checkOpt is the functional-option type for CheckErr.
type checkOpt func(*checkParams)

// checkParams holds all configurable values for CheckErr.
type checkParams struct {
	op        string
	category  string
	message   string
	details   map[string]any
	writer    io.Writer
	exitCode  int
	formatter FormatterFunc
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// WithOp lets you override the operation name that CheckErr will wrap with.
func WithOp(opName string) checkOpt {
	return func(p *checkParams) {
		p.op = opName
	}
}

// WithCategory overrides the error category.
func WithCategory(cat string) checkOpt {
	return func(p *checkParams) {
		p.category = cat
	}
}

// WithMessage overrides the user‐facing message.
func WithMessage(msg string) checkOpt {
	return func(p *checkParams) {
		p.message = msg
	}
}

// WithDetails replaces the metadata map. If you want to merge instead of
// replace, read your existing details with Details(err) first.
func WithDetails(d map[string]any) checkOpt {
	return func(p *checkParams) {
		p.details = d
	}
}

// WithWriter redirects the error output (defaults to stderr).
func WithWriter(w io.Writer) checkOpt {
	return func(p *checkParams) {
		p.writer = w
	}
}

// WithExitCode sets a custom exit code (defaults to 1).
func WithExitCode(code int) checkOpt {
	return func(p *checkParams) {
		p.exitCode = code
	}
}

// WithFormatter lets you choose any FormatterFunc (JSONFormatter, PlainFormatter,
// your own custom FormatterFunc, etc). Defaults to PseudoJSONFormatter.
func WithFormatter(f FormatterFunc) checkOpt {
	return func(p *checkParams) {
		p.formatter = f
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckErr registers, wraps, formats and logs a fatal error via Horus.
// If err is non-nil it prints using the configured FormatterFunc, then exits.
func CheckErr(err error, opts ...checkOpt) {
	if err == nil {
		return
	}

	// 1) metrics / instrumentation
	RegisterError(err)

	// 2) default parameters
	cfg := checkParams{
		op:        "check error",
		category:  "runtime_error",
		message:   "An error occurred during execution",
		details:   map[string]any{"severity": "critical", "location": "checkErr"},
		writer:    os.Stderr,
		exitCode:  1,
		formatter: PseudoJSONFormatter,
	}

	// 3) apply user overrides
	for _, opt := range opts {
		opt(&cfg)
	}

	// 4) build a rich *Herror
	herr := NewCategorizedHerror(
		cfg.op,
		cfg.category,
		cfg.message,
		err,
		cfg.details,
	)

	// 5) format & print
	if he, ok := AsHerror(herr); ok {
		fmt.Fprintln(cfg.writer, cfg.formatter(he))
	} else {
		// shouldn't happen, but fallback to plain Error()
		fmt.Fprintln(cfg.writer, herr.Error())
	}

	// 6) exit
	exitFunc(cfg.exitCode)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
