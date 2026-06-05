package harnessobs

import (
	"net/url"
	"strings"
)

// unsafeURL rejects authenticated URLs while allowing ordinary URLs that do
// not carry credentials or auth-like query keys.
func unsafeURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" {
		return false
	}
	if parsed.User != nil {
		return true
	}
	return queryHasAuthKey(parsed.Query())
}

// queryHasAuthKey treats auth-like query parameter names as credentials even
// when the URL has no userinfo section.
func queryHasAuthKey(values url.Values) bool {
	for key := range values {
		if authQueryKeys[strings.ToLower(key)] {
			return true
		}
	}
	return false
}
