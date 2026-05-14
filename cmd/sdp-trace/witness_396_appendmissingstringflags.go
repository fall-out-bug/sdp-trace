package main

import (
	"strings"
)

func appendMissingStringFlags(missing []string, opts *flagSet, required map[string]string) []string {
	for name, flag := range required {
		if strings.TrimSpace(opts.stringValue(name)) == "" {
			// Preserve the user-facing flag spelling in remediation output.
			missing = append(missing, flag)
		}
	}
	return missing
}
