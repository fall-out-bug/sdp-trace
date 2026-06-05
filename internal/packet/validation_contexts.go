package packet

import "time"

type demoFirstPacketChecker struct {
	bundle     Bundle
	now        time.Time
	rows       map[string]Row
	entryByRef map[string]BundleEntry
	errors     []string
}

type bundleValidator struct {
	bundle        Bundle
	now           time.Time
	entryByRef    map[string]BundleEntry
	resolverByRef map[string]string
	rows          map[string]Row
	errors        []string
}
