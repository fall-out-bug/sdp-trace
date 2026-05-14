package main

import (
	"fmt"
)

func previewActionsForInputs(inputs map[string]string, order []string, missingMessage, invalidMessage string) []string {
	var actions []string
	for _, key := range order {
		// Fixed key order keeps preview remediation stable for docs/tests.
		switch inputs[key] {
		case "absent":
			// Missing setup evidence is actionable before any assessment runs.
			actions = append(actions, fmt.Sprintf(missingMessage, key))
		case "present_unreadable", "present_malformed":
			// Unreadable or malformed setup evidence must be replaced, not
			// interpreted as a negative assessment condition.
			actions = append(actions, fmt.Sprintf(invalidMessage, key))
		}
	}
	return actions
}
