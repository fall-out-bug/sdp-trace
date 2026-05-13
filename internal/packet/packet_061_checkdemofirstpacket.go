package packet

import (
	"time"
)

func CheckDemoFirstPacket(bundle Bundle, now time.Time) Validation {
	// CheckDemoFirstPacket keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	validation := Validate(bundle, now)
	check := demoFirstPacketChecker{
		bundle:     bundle,
		now:        now,
		rows:       map[string]Row{},
		entryByRef: map[string]BundleEntry{},
		errors:     append([]string(nil), validation.Errors...),
	}
	return check.validate()
}
