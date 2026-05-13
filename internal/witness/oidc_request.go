package witness

import (
	"fmt"
	"net/http"
	"net/url"
)

func buildOIDCTokenRequest(requestURL, requestToken string) (*http.Request, error) {
	// The GitHub-provided request token is only sent to the expected Actions OIDC
	// host after audience normalization.
	parsed, err := url.Parse(requestURL)
	if err != nil {
		return nil, err
	}
	if err := validateOIDCTokenEndpoint(parsed); err != nil {
		return nil, err
	}
	setOIDCTokenAudience(parsed)

	req, err := http.NewRequest(http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "bearer "+requestToken)
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func validateOIDCTokenEndpoint(parsed *url.URL) error {
	if parsed.Scheme != "https" {
		// The request token must not cross plaintext transport.
		return fmt.Errorf("unexpected oidc request scheme: %s", parsed.Scheme)
	}
	if parsed.Hostname() != "token.actions.githubusercontent.com" {
		// Only GitHub's exact OIDC endpoint is allowed to receive the request token.
		return fmt.Errorf("unexpected oidc request host: %s", parsed.Host)
	}
	return nil
}

func setOIDCTokenAudience(parsed *url.URL) {
	// Audience pinning prevents replaying a token minted for another relying
	// party into this witness record.
	query := parsed.Query()
	query.Set("audience", "sdp-trace")
	parsed.RawQuery = query.Encode()
}
