package authority

import (
	"path"
	"regexp"
	"strings"
)

func targetMatches(pattern, target string) bool {
	// targetMatches keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if pattern == "" || target == "" {

		return false
	}
	if !strings.Contains(pattern, "**") {

		return pathMatches(pattern, target)
	}
	return recursivePathMatches(pattern, target)
}

func pathMatches(pattern, target string) bool {
	ok, err := path.Match(pattern, target)
	return err == nil && ok
}

func recursivePathMatches(pattern, target string) bool {
	// recursivePathMatches keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	re := regexp.QuoteMeta(pattern)
	re = strings.ReplaceAll(re, `\*\*`, `.*`)
	re = strings.ReplaceAll(re, `\*`, `[^/]*`)
	ok, err := regexp.MatchString("^"+re+"$", target)
	return err == nil && ok
}
