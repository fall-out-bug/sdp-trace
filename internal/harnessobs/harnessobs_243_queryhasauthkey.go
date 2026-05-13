package harnessobs

import (
	"net/url"

	"strings"
)

func queryHasAuthKey(values url.Values) bool {
	// queryHasAuthKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for key := range values {

		if authQueryKeys[strings.ToLower(key)] {
			return true
		}
	}
	return false
}
