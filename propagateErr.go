////////////////////////////////////////////////////////////////////////////////////////////////////

package horus

////////////////////////////////////////////////////////////////////////////////////////////////////

// PropagateErr wraps a non-nil error in an Herror with the given context.
// If err is already an Herror, its Category and Details are optionally
// preserved (unless overridden) and merged with the new details.
// If err is nil, PropagateErr returns nil.
func PropagateErr(
	op, category, message string,
	err error,
	details map[string]any,
) error {
	if err == nil {
		return nil
	}

	// Determine base Category and Details if err is already an Herror
	var baseCat string
	var baseDetails map[string]any
	if herr, ok := AsHerror(err); ok {
		baseCat = herr.Category
		baseDetails = herr.Details
	}

	// Override category if provided
	if category != "" {
		baseCat = category
	}

	// Merge details: copy baseDetails then overlay new details
	merged := make(map[string]any, len(baseDetails)+len(details))
	for k, v := range baseDetails {
		merged[k] = v
	}
	for k, v := range details {
		merged[k] = v
	}

	// Use the internal constructor so we always get a *Herror with a stack trace
	return newHerror(op, baseCat, message, err, merged)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
