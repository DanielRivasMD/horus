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

	// Collect fields excluding Stack
	fields := []field{
		{"Op", fmt.Sprintf("\"%s\"", h.Op), chalk.Yellow},
		{"Message", fmt.Sprintf("\"%s\"", h.Message), chalk.Yellow},
		{"Err", fmt.Sprintf("%v", h.Err), chalk.Yellow},
		{"Category", fmt.Sprintf("\"%s\"", h.Category), chalk.Yellow},
	}

	// Include detail fields as pseudo top-level
	var detailFields []field
	for k, v := range h.Details {
		detailFields = append(detailFields, field{
			key:   k,
			value: fmt.Sprintf("\"%v\"", v),
			color: chalk.White,
		})
	}

	// Compute max width of plain (non-colored) key names
	maxLen := 0
	for _, f := range append(fields, detailFields...) {
		if len(f.key) > maxLen {
			maxLen = len(f.key)
		}
	}

	// Print aligned top-level fields
	for _, f := range fields[:3] {
		paddedKey := fmt.Sprintf("%-*s", maxLen, f.key)
		fmt.Fprintf(&b, "%s %s,\n", f.color.Color(paddedKey), chalk.Red.Color(f.value))
	}

	// Print Details
	b.WriteString(chalk.Yellow.Color("Details") + "\n")
	for _, f := range detailFields {
		paddedKey := fmt.Sprintf("  %-*s", maxLen, f.key)
		fmt.Fprintf(&b, "%s %s,\n", f.color.Color(paddedKey), chalk.Red.Color(f.value))
	}

	// Print Category
	paddedKey := fmt.Sprintf("%-*s", maxLen, fields[3].key)
	fmt.Fprintf(&b, "%s %s,\n", fields[3].color.Color(paddedKey), chalk.Red.Color(fields[3].value))

	// Stack remains unaligned
	b.WriteString(chalk.Yellow.Color("Stack") + "\n")
	for _, addr := range h.Stack {
		b.WriteString("  " + chalk.Dim.TextStyle(fmt.Sprintf("%v", addr)) + "\n")
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
