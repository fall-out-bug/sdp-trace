package repoobserver

import "strings"

func sanitizeOrigin(origin string) string {
	// URL credentials are stripped before repository identity hashing.
	origin = removeOriginFragment(strings.TrimSpace(origin))
	hadURLCredentials := hasURLCredentials(origin)
	origin = removeOriginCredentials(origin)
	if hadURLCredentials {
		return origin
	}
	return originTail(origin)
}
