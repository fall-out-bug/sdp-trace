package packet

import (
	"strings"
	"time"
)

func demoUsableEntry(entry BundleEntry, now time.Time) bool {
	return entryHasResolverAndDigest(entry) && !entryExpired(entry, now) && !passRefUnverifiable(entry)
}

func entryHasResolverAndDigest(entry BundleEntry) bool {
	return strings.TrimSpace(entry.Resolver) != "" && strings.TrimSpace(entry.Digest) != ""
}

func syntheticEntryDigest(entry BundleEntry) bool {
	return strings.TrimSpace(entry.Digest) == "" || entry.Digest == digestPlaceholder(entry.Ref+entry.Resolver)
}
