package repoobserver

import "strings"

func removeOriginFragment(origin string) string {
	if idx := strings.Index(origin, "#"); idx >= 0 {
		// URL fragments are local navigation hints and not repository identity.
		return origin[:idx]
	}
	return origin
}

func removeOriginCredentials(origin string) string {
	// SCP-like and URL remotes encode userinfo differently, so redact both forms.
	if strings.Contains(origin, "@") && !strings.Contains(origin, "://") {
		return origin[strings.LastIndex(origin, "@")+1:]
	}
	if originHasURLCredentials(origin) {
		return originWithoutURLCredentials(origin)
	}
	return origin
}

func originHasURLCredentials(origin string) bool {
	// Treat @ as URL credentials only when it appears after a URL scheme.
	at := strings.LastIndex(origin, "@")
	return at >= 0 && strings.Contains(origin[:max(at, 0)], "://")
}
