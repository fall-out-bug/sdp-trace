package main

import (
	"net/url"
	"strings"
)

func localHTTPGitHubAPI(parsed *url.URL) bool {
	return parsed.Scheme == "http" && loopbackHost(strings.ToLower(parsed.Hostname()))
}
