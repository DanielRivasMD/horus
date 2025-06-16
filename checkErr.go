////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"os"
)

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

// CheckErr reports errors using the horus library and exits gracefully.
func CheckErr(err error) {
	if err != nil {
		// Register the error type.
		RegisterError(err)
		herr := NewCategorizedHerror(
			"check error",
			"runtime_error",
			"An error occurred during execution",
			err,
			map[string]any{"severity": "critical", "location": "checkErr function"},
		)

		// Log the error in a colored format for easier debugging.
		fmt.Println(FormatError(herr, SimpleColoredFormatter))
		os.Exit(1) // Exit gracefully
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
