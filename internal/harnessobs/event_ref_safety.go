package harnessobs

import "strings"

func safeEventRef(ref string) bool {
	if unsafeEventRefPath(ref) {
		return false
	}
	if !strings.HasPrefix(ref, "events/") || !strings.HasSuffix(ref, ".json") {
		return false
	}

	id := strings.TrimSuffix(strings.TrimPrefix(ref, "events/"), ".json")
	return safeFileIDPattern.MatchString(id)
}
