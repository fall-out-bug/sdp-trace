package witness

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func buildOIDCTokenRequest(requestURL, requestToken string) (*http.Request, error) {
	// The GitHub-provided request token is only sent to the expected Actions OIDC
	// host after audience normalization.
	parsed, err := url.Parse(requestURL)
	if err != nil {
		return nil, err
	}
	if err := validateOIDCTokenHost(parsed.Host); err != nil {
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

func validateOIDCTokenHost(host string) error {
	if !strings.HasSuffix(host, "actions.githubusercontent.com") {
		// Only GitHub's OIDC endpoint is allowed to receive the request token.
		return fmt.Errorf("unexpected oidc request host: %s", host)
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
