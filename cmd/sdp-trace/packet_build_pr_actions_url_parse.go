package main

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

func requireHTTPSGitHubAPI(parsed *url.URL, apiURL string) error {
	if parsed.Scheme == "https" {
		// HTTPS is the only scheme allowed for credential-bearing GitHub calls.
		return nil
	}
	return fmt.Errorf("unsafe GitHub API URL %q: HTTPS is required before sending GitHub credentials", apiURL)
}

func parseGitHubAPIURL(apiURL string) (*url.URL, error) {
	parsed, err := url.Parse(apiURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("unsafe GitHub API URL %q", apiURL)
	}
	// The parsed URL is still untrusted until scheme, credentials, and host are
	// validated by the caller.
	return parsed, nil
}

func localHTTPGitHubAPI(parsed *url.URL) bool {
	return parsed.Scheme == "http" && loopbackHost(strings.ToLower(parsed.Hostname()))
}

func loopbackHost(host string) bool {
	ip := net.ParseIP(host)
	return host == "localhost" || (ip != nil && ip.IsLoopback())
}
