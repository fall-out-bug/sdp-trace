package packet

import (
	"time"
)

func demoUsableEntry(entry BundleEntry, now time.Time) bool {

	return entryHasResolverAndDigest(entry) && !entryExpired(entry, now) && !passRefUnverifiable(entry)
}
