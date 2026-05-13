package main

import (
	"fmt"
)

func protectedPreviewActions(inputs map[string]string) []string {
	names := []string{"checkpoint", "checkpoint_policy", "witness"}
	actions := make([]string, 0)
	for _, name := range names {
		// Fixed input order keeps preview remediation stable for docs/tests.
		switch inputs[name] {
		case "absent":
			actions = append(actions, fmt.Sprintf("Supply %s input before running protected gate.", name))
		case "present_unreadable", "present_malformed":
			// Both unreadable and malformed artifacts require replacement before
			// protected replay can make an evidence-backed claim.
			actions = append(actions, fmt.Sprintf("Replace %s input with readable JSON.", name))
		}
	}
	return actions
}
