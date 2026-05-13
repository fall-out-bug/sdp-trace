package main

import (
	"errors"
	"net/url"
)

func validateParsedGitHubAPIURL(parsed *url.URL, apiURL, serverURL string) error {
	if parsed.User != nil {
		// Credentials must travel through Authorization headers, never URLs that
		// can leak through logs, errors, or packet context.
		return errors.New("unsafe GitHub API URL: embedded credentials are not allowed")
	}
	return validateGitHubAPIURLTrustTarget(parsed, apiURL, serverURL)
}
