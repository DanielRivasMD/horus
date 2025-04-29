////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckFileExist checks if a file exists. If it doesn't, it executes
// the provided notFoundAction. Any error from the stat operation or the
// notFoundAction will be wrapped in an Herror.
func CheckFileExist(filePath string, notFoundAction NotFoundAction) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File '%s' does not exist. Executing custom action.\n", filePath)
			if notFoundAction != nil {
				if actionErr := notFoundAction(filePath); actionErr != nil {
					return NewHerror("check file", fmt.Sprintf("'%s' not found action failed", filePath), actionErr, map[string]any{"path": filePath})
				}
				return nil // Action succeeded, no error to report for the check itself
			}
			return NewHerror("check file", fmt.Sprintf("file '%s' not found", filePath), nil, map[string]any{"path": filePath, "action": "none"})
		}
		return NewHerror("check file", fmt.Sprintf("error checking file '%s'", filePath), err, map[string]any{"path": filePath})
	}
	fmt.Printf("File '%s' exists: '%s'\n", filePath, filePath)
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
