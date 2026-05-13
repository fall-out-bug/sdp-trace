package witness

import (
	"sort"
	"strings"
)

func missingGitHubOIDC(env map[string]string) []string {
	// Both request fields are required before a live GitHub OIDC token can be
	// fetched and checked.
	required := []string{"ACTIONS_ID_TOKEN_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"}
	missing := make([]string, 0)
	for _, key := range required {
		// Treat blank token request values as absent evidence.
		if strings.TrimSpace(env[key]) == "" {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)
	return missing
}
