package main

import (
	"strings"
)

func missingEnvFields(env map[string]string, required []string) []string {
	missing := make([]string, 0)
	for _, key := range required {
		if strings.TrimSpace(env[key]) == "" {
			// Trim whitespace so empty exported variables are treated as absent.
			missing = append(missing, key)
		}
	}
	return missing
}
