package witness

import "strings"

func missingEnvKeys(env map[string]string, required []string) []string {
	// Required identity fields are validated by name so missing evidence can be
	// reported directly to callers and trace records.
	missing := make([]string, 0)
	for _, key := range required {
		// Whitespace-only CI fields cannot bind the source or run identity.
		if strings.TrimSpace(env[key]) == "" {
			missing = append(missing, key)
		}
	}
	return missing
}
