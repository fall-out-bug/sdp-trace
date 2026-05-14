package main

func previewActionForInputState(inputs map[string]string, key, absentAction, malformedAction string) []string {
	// Preview action helpers translate setup state into remediation text only;
	// they do not evaluate the underlying assessment payload.
	switch inputs[key] {
	case "absent":
		return []string{absentAction}
	case "present_unreadable", "present_malformed":
		// Unreadable and malformed inputs share the same repair path because
		// both block assessment before verdict logic can run.
		return []string{malformedAction}
	default:
		return nil
	}
}
