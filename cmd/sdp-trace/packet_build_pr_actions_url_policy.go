package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func validateGitHubAPIURL(apiURL, serverURL string) error {
	parsed, err := parseGitHubAPIURL(apiURL)
	if err != nil {
		return err
	}
	// Parse and target validation are split so error text can distinguish syntax
	// from trust-target failures.
	return validateParsedGitHubAPIURL(parsed, apiURL, serverURL)
}

func validateParsedGitHubAPIURL(parsed *url.URL, apiURL, serverURL string) error {
	if parsed.User != nil {
		// Credentials must travel through Authorization headers, never URLs that
		// can leak through logs, errors, or packet context.
		return errors.New("unsafe GitHub API URL: embedded credentials are not allowed")
	}
	return validateGitHubAPIURLTrustTarget(parsed, apiURL, serverURL)
}

func validateGitHubAPIURLTrustTarget(parsed *url.URL, apiURL, serverURL string) error {
	if localHTTPGitHubAPI(parsed) {
		// Loopback HTTP is allowed for hermetic tests; non-local API targets must
		// satisfy HTTPS and host binding before receiving credentials.
		return nil
	}
	if err := requireHTTPSGitHubAPI(parsed, apiURL); err != nil {
		return err
	}
	if githubAPIHostAllowed(strings.ToLower(parsed.Hostname()), serverURL) {
		// Allowed hosts are either public GitHub's API host or the configured
		// Enterprise host.
		return nil
	}
	return fmt.Errorf("unsafe GitHub API URL %q: host is not the configured GitHub host", apiURL)
}
