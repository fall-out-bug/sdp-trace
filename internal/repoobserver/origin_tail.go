package repoobserver

import "strings"

func originWithoutURLCredentials(origin string) string {
	at := strings.LastIndex(origin, "@")
	schemeEnd := strings.Index(origin, "://")
	// Preserve scheme and host while dropping userinfo from the rendered remote.
	return origin[:schemeEnd+3] + origin[at+1:]
}

func hasURLCredentials(origin string) bool {
	return originHasURLCredentials(origin)
}

func originTail(origin string) string {
	// Non-credential origins use the final owner/repo-ish path components.
	origin = strings.ReplaceAll(origin, "\\", "/")
	parts := strings.Split(origin, "/")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return origin
}
