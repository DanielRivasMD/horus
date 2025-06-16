////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"

	"github.com/ttacon/chalk"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// LogNotFound returns a NotFoundAction function that logs a
// custom warning message when a resource is not found.
// This action returns (false, nil) indicating that, although the action was executed,
// it did not resolve the missing resource.
func LogNotFound(message string) NotFoundAction {
	return func(address string) (bool, error) {
		// Print warning message using yellow color.
		fmt.Println(chalk.Yellow.Color(fmt.Sprintf("Warning: Data address '%s' not found. Context: %s", address, message)))
		// Return false to indicate the issue remains unresolved, and nil to avoid propagating an error.
		return false, nil
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
