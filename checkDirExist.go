////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckDirExist checks if a directory exists. If it doesn't, it executes the provided
// notFoundAction. It returns a tuple: a boolean indicating overall success (either the directory exists
// or the custom action adequately handled the situation) and an error carrying detailed diagnostic context.
// The verbose flag enables optional logging.
func CheckDirExist(dirPath string, notFoundAction NotFoundAction, verbose bool) (bool, error) {
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if verbose {
				fmt.Printf("Directory '%s' does not exist. Executing custom action.\n", dirPath)
			}
			if notFoundAction != nil {
				ok, actionErr := notFoundAction(dirPath)
				if actionErr != nil {
					// Wrap the error with additional diagnostic context.
					return false, NewCategorizedHerror(
						"check directory",
						"directory_error",
						fmt.Sprintf("'%s' not found action failed", dirPath),
						actionErr,
						map[string]any{"path": dirPath},
					)
				}
				// Let the custom action decide the boolean outcome.
				return ok, nil
			}
			// Directory not found and no custom action provided.
			return false, NewCategorizedHerror(
				"check directory",
				"directory_error",
				fmt.Sprintf("directory '%s' not found", dirPath),
				nil,
				map[string]any{"path": dirPath, "action": "none"},
			)
		}
		// Some other error occurred during Stat.
		return false, NewCategorizedHerror(
			"check directory",
			"directory_error",
			fmt.Sprintf("error checking directory '%s'", dirPath),
			err,
			map[string]any{"path": dirPath},
		)
	}
	if verbose {
		fmt.Printf("Directory '%s' exists.\n", dirPath)
	}
	return true, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
