package packet

import (
	"strings"
)

func checkResolvers(checks []GitHubCheck) string {
	// checkResolvers keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	values := []string{}
	for _, check := range checks {

		values = append(values, check.Name+"="+check.URL)
	}
	return strings.Join(values, ", ")
}
