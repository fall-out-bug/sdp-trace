package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func protectedInputStatus(path string) string {
	if strings.TrimSpace(path) == "" {
		// Missing protected inputs are reported as setup gaps, not as verdicts.
		return "absent"
	}
	var value any
	if err := readJSONFile(path, &value); err != nil {
		return protectedInputErrorStatus(err)
	}
	return "present_readable"
}

func protectedInputErrorStatus(err error) string {
	if os.IsNotExist(err) || errors.Is(err, os.ErrPermission) {
		// Protected preview distinguishes unavailable inputs from malformed JSON
		// so users know whether to fix access or content.
		return "present_unreadable"
	}
	return "present_malformed"
}

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
