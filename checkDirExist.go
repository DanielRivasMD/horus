////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

// NotFoundAction is a function type that takes the directory path as input
// and returns an error.
type NotFoundAction func(string) error

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
					return NewHerror("check directory", fmt.Sprintf("'%s' not found action failed", dirPath), actionErr, map[string]interface{}{"path": dirPath})
				}
				return nil // Action succeeded, no error to report for the check itself
			}
			return NewHerror("check directory", fmt.Sprintf("directory '%s' not found", dirPath), nil, map[string]interface{}{"path": dirPath, "action": "none"})
		}
		return NewHerror("check directory", fmt.Sprintf("error checking directory '%s'", dirPath), err, map[string]interface{}{"path": dirPath})
	}
	fmt.Printf("Directory '%s' exists: '%s'\n", dirPath, dirPath)
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
