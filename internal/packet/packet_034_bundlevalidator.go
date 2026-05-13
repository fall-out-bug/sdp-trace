package packet

import (
	"time"
)

type bundleValidator struct {
	bundle        Bundle
	now           time.Time
	entryByRef    map[string]BundleEntry
	resolverByRef map[string]string
	rows          map[string]Row
	errors        []string
}
