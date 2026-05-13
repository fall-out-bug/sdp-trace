package harnessobs

import (
	"strings"
)

func safeToken(value string) string {
	// safeToken keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var b strings.Builder
	for _, r := range value {
		writeSafeTokenRune(&b, r)
		if b.Len() >= 128 {

			break
		}
	}
	token := strings.Trim(b.String(), "-_.:")
	if token == "" {

		return "opencode"
	}
	return token
}
