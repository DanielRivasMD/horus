////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// CheckDirExist checks if a directory exists. If it doesn't, it executes
// the provided notFoundAction. Any error from the stat operation or the
// notFoundAction will be wrapped in an Herror.
func CheckDirExist(dirPath string, notFoundAction NotFoundAction) error {
	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Directory '%s' does not exist. Executing custom action.\n", dirPath)
			if notFoundAction != nil {
				if actionErr := notFoundAction(dirPath); actionErr != nil {
					return NewCategorizedHerror(
						"check directory",
						"directory_error",
						fmt.Sprintf("'%s' not found action failed", dirPath),
						actionErr,
						map[string]any{"path": dirPath},
					)
				}
				return nil // Action succeeded, no error to report for the check itself
			}
			return NewCategorizedHerror(
				"check directory",
				"directory_error",
				fmt.Sprintf("directory '%s' not found", dirPath),
				nil,
				map[string]any{"path": dirPath, "action": "none"},
			)
		}
		return NewCategorizedHerror(
			"check directory",
			"directory_error",
			fmt.Sprintf("error checking directory '%s'", dirPath),
			err,
			map[string]any{"path": dirPath},
		)
	}
	fmt.Printf("Directory '%s' exists: '%s'\n", dirPath, dirPath)
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
