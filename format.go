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

	writeField := func(key string, value string) {
		fmt.Fprintf(&b, "  %s: %s,\n", chalk.White.Color(key), chalk.Red.Color(value))
	}

	b.WriteString("")

	writeField("Op", fmt.Sprintf("\"%s\"", h.Op))
	writeField("Message", fmt.Sprintf("\"%s\"", h.Message))
	writeField("Err", fmt.Sprintf("%v", h.Err))

	// Format Details map
	b.WriteString("  " + chalk.White.Color("Details") + ": {\n")
	for k, v := range h.Details {
		fmt.Fprintf(&b, "    %s: %s,\n",
			chalk.White.Color(k),
			chalk.Red.Color(fmt.Sprintf("\"%v\"", v)))
	}
	b.WriteString("  },\n")

	writeField("Category", fmt.Sprintf("\"%s\"", h.Category))

	// Format Stack
	b.WriteString("  " + chalk.White.Color("Stack") + ": [\n")
	for _, addr := range h.Stack {
		fmt.Fprintf(&b, "    %v,\n", addr)
	}
	b.WriteString("  ]\n")

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
