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
	}

	// Collect top-level fields
	fields := []field{
		{"Op", fmt.Sprintf("\"%s\"", h.Op)},
		{"Message", fmt.Sprintf("\"%s\"", h.Message)},
		{"Err", fmt.Sprintf("%v", h.Err)},
		{"Category", fmt.Sprintf("\"%s\"", h.Category)},
	}

	// Calculate maximum width for alignment
	maxWidth := 0
	for _, f := range fields {
		if len(f.key) > maxWidth {
			maxWidth = len(f.key)
		}
	}
	for k := range h.Details {
		if len(k) > maxWidth-2 { // +2 because of indent
			maxWidth = len(k) + 2
		}
	}

	// Format top-level fields
	for _, f := range fields[:3] { // Op, Message, Err
		fmt.Fprintf(&b, "%-*s%s%s,\n", maxWidth+2, chalk.Yellow.Color(f.key+":"), " ", chalk.Red.Color(f.value))
	}

	// Format Details
	b.WriteString(chalk.Yellow.Color("Details:") + "\n")
	for k, v := range h.Details {
		key := fmt.Sprintf("  %s:", k)
		fmt.Fprintf(&b, "%-*s%s%s,\n", maxWidth+2, chalk.White.Color(key), " ", chalk.Red.Color(fmt.Sprintf("\"%v\"", v)))
	}
	b.WriteString("\n")

	// Category
	fmt.Fprintf(&b, "%-*s%s%s,\n", maxWidth+2, chalk.Yellow.Color("Category:"), " ", chalk.Red.Color(fields[3].value))

	// Stack
	b.WriteString(chalk.Yellow.Color("Stack:") + "\n")
	for _, addr := range h.Stack {
		b.WriteString("          " + chalk.Red.Color(fmt.Sprintf("%v", addr)) + "\n")
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
