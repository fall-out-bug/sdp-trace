package packet

import (
	"time"
)

func Validate(bundle Bundle, now time.Time) Validation {
	// Validate keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	validator := bundleValidator{
		bundle:        bundle,
		now:           now,
		entryByRef:    map[string]BundleEntry{},
		resolverByRef: map[string]string{},
	}
	return validator.validate()
}
