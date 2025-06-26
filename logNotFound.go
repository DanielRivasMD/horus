////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"io"
	"os"

	"github.com/ttacon/chalk"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

type logNotFoundConfig struct {
	writer   io.Writer
	template string
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// LogNotFound returns a NotFoundAction that logs a warning when a resource
// isn’t found.  By default it prints to stderr in yellow:
//
//	Warning: Data address '...' not found. Context: ...
func LogNotFound(contextMsg string, opts ...LogNotFoundOption) NotFoundAction {
	cfg := logNotFoundConfig{
		writer:   os.Stderr,
		template: "Warning: Data address '%s' not found. Context: %s",
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	return func(address string) (bool, error) {
		msg := fmt.Sprintf(cfg.template, address, contextMsg)
		msg = chalk.Yellow.Color(msg)
		fmt.Fprintln(cfg.writer, msg)
		return false, nil
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// LogNotFoundOption customizes how LogNotFound prints.
type LogNotFoundOption func(*logNotFoundConfig)

////////////////////////////////////////////////////////////////////////////////////////////////////

// WithLogWriter directs the warning message to a different io.Writer.
func WithLogWriter(w io.Writer) LogNotFoundOption {
	return func(cfg *logNotFoundConfig) {
		cfg.writer = w
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// WithNotFoundTemplate lets you override the printf-style template.
func WithNotFoundTemplate(tmpl string) LogNotFoundOption {
	return func(cfg *logNotFoundConfig) {
		cfg.template = tmpl
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
