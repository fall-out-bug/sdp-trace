package packet

import (
	"strings"
)

func entryHasResolverAndDigest(entry BundleEntry) bool {
	return strings.TrimSpace(entry.Resolver) != "" && strings.TrimSpace(entry.Digest) != ""
}
