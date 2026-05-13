package packet

import (
	"strings"
	"time"
)

func entryExpired(entry BundleEntry, now time.Time) bool {
	// entryExpired keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if strings.TrimSpace(entry.ExpiresAt) == "" {
		return false
	}
	expiresAt, err := time.Parse(time.RFC3339, entry.ExpiresAt)
	if err != nil {

		return true
	}
	return !expiresAt.After(now)
}
