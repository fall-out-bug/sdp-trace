package authority

import (
	"sort"
	"strings"
)

func uniqueStrings(values []string) []string {
	// uniqueStrings keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {

			continue
		}
		seen[value] = true
		out = append(out, value)
	}

	sort.Strings(out)
	return out
}

func mapKeys(values map[string]bool) []string {
	// mapKeys keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	out := make([]string, 0, len(values))
	for key := range values {

		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
