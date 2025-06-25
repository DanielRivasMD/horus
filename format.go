////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ttacon/chalk"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// FormatterFunc defines a function type for custom error formatting.
type FormatterFunc func(*Herror) string

////////////////////////////////////////////////////////////////////////////////////////////////////

// FormatError formats the Herror using a custom formatter function.
func FormatError(err error, formatter FormatterFunc) string {
	if herr, ok := AsHerror(err); ok {
		return formatter(herr)
	}
	return err.Error() // Fallback to default error string if not an Herror
}

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

	// Collect all top-level fields (excluding Stack)
	fields := []field{
		{"Op", fmt.Sprintf("\"%s\"", h.Op), chalk.Yellow},
		{"Message", fmt.Sprintf("\"%s\"", h.Message), chalk.Yellow},
		{"Err", fmt.Sprintf("%v", h.Err), chalk.Yellow},
		{"Category", fmt.Sprintf("\"%s\"", h.Category), chalk.Yellow},
	}

	// Flatten Details map into same format
	var detailFields []field
	for k, v := range h.Details {
		detailFields = append(detailFields, field{
			key:   k,
			value: fmt.Sprintf("\"%v\"", v),
			color: chalk.White,
		})
	}

	// Find longest key for padding
	maxLen := 0
	for _, f := range append(fields, detailFields...) {
		if len(f.key) > maxLen {
			maxLen = len(f.key)
		}
	}

	// Format top-level fields
	for _, f := range fields[:3] { // Op, Message, Err
		fmt.Fprintf(&b, "%-*s %s,\n", maxLen+1, f.color.Color(f.key+":"), chalk.Red.Color(f.value))
	}

	// Format Details
	b.WriteString(chalk.Yellow.Color("Details:") + "\n")
	for _, f := range detailFields {
		fmt.Fprintf(&b, "  %-*s %s,\n", maxLen+1, f.color.Color(f.key+":"), chalk.Red.Color(f.value))
	}
	b.WriteString("\n")

	// Category
	fmt.Fprintf(&b, "%-*s %s,\n", maxLen+1, fields[3].color.Color(fields[3].key+":"), chalk.Red.Color(fields[3].value))

	// Format Stack (unaligned, indented)
	b.WriteString(chalk.Yellow.Color("Stack:") + "\n")
	for _, addr := range h.Stack {
		b.WriteString("  " + chalk.Red.Color(fmt.Sprintf("%v", addr)) + "\n")
	}
	b.WriteString("\n")

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
