package packet

import (
	"net/url"
	"strings"
)

// The redaction check is deliberately small and conservative; it prevents
// obvious secret-like resolver strings from being persisted in generated
// examples or packets.
func redactSecretLike(value string) string {
	if containsSecretMarker(value) || containsURLCredentials(value) {
		return "[redacted-secret]"
	}
	return value
}

func containsSecretMarker(value string) bool {
	upper := strings.ToUpper(value)
	for _, marker := range secretResolverMarkers() {
		if strings.Contains(upper, marker) {
			return true
		}
	}
	return false
}

func secretResolverMarkers() []string {
	return []string{"SECRET", "TOKEN", "AUTHORIZATION:", "ACCESS_TOKEN=", "PASSWORD=", "SECRET=", "TOKEN=", "AUTH="}
}

func containsURLCredentials(value string) bool {
	parsed, err := url.Parse(value)
	if err == nil && parsed.User != nil {
		return true
	}
	return false
}
