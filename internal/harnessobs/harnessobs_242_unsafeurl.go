package harnessobs

import (
	"net/url"
)

func unsafeURL(raw string) bool {
	// unsafeURL keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" {
		return false
	}
	if parsed.User != nil {

		return true
	}
	return queryHasAuthKey(parsed.Query())
}
