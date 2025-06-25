////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/json"
	"fmt"

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

func PrettyColoredJSONFormatter(h *Herror) string {
	jsonBytes, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return chalk.Red.Color(fmt.Sprintf("Error formatting JSON: %v", err))
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &raw); err != nil {
		return chalk.Red.Color(fmt.Sprintf("Error parsing JSON: %v", err))
	}

	var result string
	result += "{\n"
	for key, value := range raw {
		fieldName := chalk.White.Color(fmt.Sprintf("  \"%s\"", key))
		var fieldValue string

		switch v := value.(type) {
		case string:
			fieldValue = chalk.Red.Color(fmt.Sprintf("\"%s\"", v))
		default:
			marshaledVal, _ := json.MarshalIndent(v, "  ", "  ")
			fieldValue = chalk.Red.Color(string(marshaledVal))
		}
		result += fmt.Sprintf("%s: %s,\n", fieldName, fieldValue)
	}
	result = result[:len(result)-2] + "\n}" // Remove last comma and close brace
	return result
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
