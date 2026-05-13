package packet

import (
	"time"
)

type demoFirstPacketChecker struct {
	bundle     Bundle
	now        time.Time
	rows       map[string]Row
	entryByRef map[string]BundleEntry
	errors     []string
}
