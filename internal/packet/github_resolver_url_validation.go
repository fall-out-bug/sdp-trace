package packet

import (
	"fmt"
	"net/url"
	"strings"
)

// Resolver URL validation is intentionally narrower than general URL parsing.
// Empty values and non-URL strings remain valid resolver labels. URL-shaped
// values must be parseable public HTTPS locations without embedded credentials,
// because generated packet evidence can be copied into review artifacts and
// should not encode localhost, private network, or credential-bearing links.

func validateResolverURL(field, value string) error {
	raw := strings.TrimSpace(value)
	if !isURLLikeResolver(raw) {
		return nil
	}
	parsed, err := parseResolverURL(raw)
	if err != nil {
		return fmt.Errorf("unsafe %s %q", field, value)
	}
	return validateParsedResolverURL(field, value, parsed)
}

func isURLLikeResolver(raw string) bool {
	return raw != "" && strings.Contains(raw, "://")
}

func parseResolverURL(raw string) (*url.URL, error) {
	// Both scheme and host are required before applying resolver policy. A value
	// like "https:///" is URL-shaped input, but it is not a portable evidence
	// reference and should fail before host policy sees an empty hostname.
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid resolver url")
	}
	return parsed, nil
}

func validateParsedResolverURL(field, value string, parsed *url.URL) error {
	// Credentials and non-HTTPS schemes are rejected before host policy so error
	// messages explain the first actionable correction. Host policy then blocks
	// local-only resolver targets that would make evidence non-portable.
	if resolverHasCredentials(parsed) {
		return fmt.Errorf("unsafe %s %q: embedded credentials are not allowed", field, value)
	}
	if parsed.Scheme != "https" {
		return fmt.Errorf("unsafe %s %q: HTTPS is required", field, value)
	}
	if unsafeResolverHost(parsed.Hostname()) {
		return fmt.Errorf("unsafe %s %q: internal hosts are not allowed", field, value)
	}
	return nil
}

func resolverHasCredentials(parsed *url.URL) bool {
	return parsed.User != nil
}
