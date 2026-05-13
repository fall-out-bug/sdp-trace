package packet

import (
	"strings"
)

func redactSecretLike(value string) string {
	// redactSecretLike keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	redacted := value
	for _, marker := range []string{"SECRET", "TOKEN", "Authorization:"} {

		if strings.Contains(strings.ToUpper(redacted), strings.ToUpper(marker)) {
			return "[redacted-secret]"
		}
	}
	return redacted
}
