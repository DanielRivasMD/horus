////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

// PropagateErr checks if err is non-nil.
// If so, it wraps the error with contextual details using NewCategorizedHerror and returns it.
// Otherwise, it returns nil.
func PropagateErr(operation, category, message string, err error, details map[string]any) error {
	if err != nil {
		return NewCategorizedHerror(operation, category, message, err, details)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
