package interaction

import (
	"regexp"
)

func unsafeCount(body []byte) int {
	// unsafeCount keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	text := string(body)
	count := 0
	for _, pattern := range []*regexp.Regexp{privatePathPattern, tokenPattern, authURLPattern} {
		count += len(pattern.FindAllString(text, -1))
	}
	return count
}
