package main

import (
	"fmt"
	"net/url"
	"strings"
)

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
