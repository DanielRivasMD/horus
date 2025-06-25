////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/ttacon/chalk"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// FormatterFunc defines a function type for custom error formatting.
type FormatterFunc func(*Herror) string

////////////////////////////////////////////////////////////////////////////////////////////////////

// JSONFormatter generates a JSON representation of an Herror.
func JSONFormatter(h *Herror) string {
	jsonOutput, err := json.Marshal(h)
	if err != nil {
		return fmt.Sprintf("error formatting: %v", err)
	}
	return string(jsonOutput)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func PseudoJSONFormatter(h *Herror) string {
	var b strings.Builder

	type field struct {
		key   string
		value string
		color chalk.Color
	}

	// Collect top-level fields (excluding Stack)
	fields := []field{
		{"Op", h.Op, chalk.Yellow},
		{"Message", h.Message, chalk.Yellow},
		{"Err", fmt.Sprintf("%v", h.Err), chalk.Yellow},
		{"Category", h.Category, chalk.Yellow},
	}

	// Convert Details into aligned field list
	var detailFields []field
	for k, v := range h.Details {
		detailFields = append(detailFields, field{
			key:   k,
			value: fmt.Sprintf("%v", v),
			color: chalk.White,
		})
	}

	// Compute max key width for padding
	maxLen := 0
	for _, f := range append(fields, detailFields...) {
		if len(f.key) > maxLen {
			maxLen = len(f.key)
		}
	}

	// Render top-level fields (except Stack)
	for _, f := range fields[:3] {
		padded := fmt.Sprintf("%-*s", maxLen, f.key)
		fmt.Fprintf(&b, "%s %s,\n", f.color.Color(padded), chalk.Red.Color(f.value))
	}

	// Render Details
	b.WriteString(chalk.Yellow.Color("Details") + "\n")
	for _, f := range detailFields {
		padded := fmt.Sprintf("  %-*s", maxLen, f.key)
		fmt.Fprintf(&b, "%s %s,\n", f.color.Color(padded), chalk.Red.Color(f.value))
	}
	b.WriteString("\n")

	// Render Category
	padded := fmt.Sprintf("%-*s", maxLen, fields[3].key)
	fmt.Fprintf(&b, "%s %s,\n", fields[3].color.Color(padded), chalk.Red.Color(fields[3].value))

	// Render Stack (show function in magenta, location dimmed)
	b.WriteString(chalk.Yellow.Color("Stack") + "\n")

	frames := runtime.CallersFrames(h.Stack)
	for {
		frame, more := frames.Next()

		// colorize parts separately
		fn := chalk.Magenta.Color(frame.Function + "()")
		loc := chalk.Dim.TextStyle(fmt.Sprintf(" %s:%d", frame.File, frame.Line))

		b.WriteString("  " + fn + loc + "\n")

		if !more {
			break
		}
	}

	return b.String()
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// SimpleColoredFormatter generates a colored representation of an Herror using the chalk library.
// You can extend this function to use different colors based on the error category.
func SimpleColoredFormatter(h *Herror) string {
	return chalk.Red.Color(fmt.Sprintf("ERROR: %s", h.Error()))
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// FormatPanic returns a red-formatted panic message for the given operation and message.
func FormatPanic(op, message string) string {
	return chalk.Red.Color(fmt.Sprintf("Panic [%s]: %s", op, message))
}

////////////////////////////////////////////////////////////////////////////////////////////////////
