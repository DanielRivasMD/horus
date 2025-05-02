////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckFileExist checks if a file exists. If it doesn't, it executes the provided
// notFoundAction. It returns a tuple: a boolean indicating whether the check is successful
// (either the file exists or the custom action adequately handled the missing file)
// and an error carrying detailed diagnostic context (if any). The verbose flag enables optional logging.
func CheckFileExist(filePath string, notFoundAction NotFoundAction, verbose bool) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			if verbose {
				fmt.Printf("File '%s' does not exist. Executing custom action.\n", filePath)
			}
			if notFoundAction != nil {
				ok, actionErr := notFoundAction(filePath)
				if actionErr != nil {
					return false, NewHerror("check file", fmt.Sprintf("'%s' not found action failed", filePath), actionErr, map[string]any{"path": filePath})
				}
				return ok, nil
			}
			return false, NewHerror("check file", fmt.Sprintf("file '%s' not found", filePath), nil, map[string]any{"path": filePath, "action": "none"})
		}
		return false, NewHerror("check file", fmt.Sprintf("error checking file '%s'", filePath), err, map[string]any{"path": filePath})
	}
	if verbose {
		fmt.Printf("File '%s' exists.\n", filePath)
	}
	return true, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
