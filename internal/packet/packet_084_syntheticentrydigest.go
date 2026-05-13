package packet

import (
	"strings"
)

func syntheticEntryDigest(entry BundleEntry) bool {
	return strings.TrimSpace(entry.Digest) == "" || entry.Digest == digestPlaceholder(entry.Ref+entry.Resolver)
}
