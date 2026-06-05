package packet

import (
	"net"
	"strings"
)

// Resolver hosts are evidence references, not fetch targets for this package.
// Blocking local and private names here prevents generated evidence from
// normalizing links that only work on the author's machine or that point at
// internal services. Public hostnames remain allowed because packet generation
// only records the resolver string.

func unsafeResolverHost(host string) bool {
	normalized := strings.ToLower(strings.TrimSpace(host))
	if unsafeResolverHostName(normalized) {
		return true
	}
	ip := net.ParseIP(normalized)
	return unsafeResolverIP(ip)
}

func unsafeResolverHostName(normalized string) bool {
	if normalized == "" || normalized == "localhost" {
		return true
	}
	return strings.HasSuffix(normalized, ".localhost") || strings.HasSuffix(normalized, ".local")
}
