package main

import (
	"fmt"
	"net/url"
)

func requireHTTPSGitHubAPI(parsed *url.URL, apiURL string) error {
	if parsed.Scheme == "https" {
		// HTTPS is the only scheme allowed for credential-bearing GitHub calls.
		return nil
	}
	return fmt.Errorf("unsafe GitHub API URL %q: HTTPS is required before sending GitHub credentials", apiURL)
}
