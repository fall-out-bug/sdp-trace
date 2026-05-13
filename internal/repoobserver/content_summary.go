package repoobserver

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func contentSummary(data []byte) string {
	// Force summaries expose digest/size metadata instead of raw file content.
	sum := sha256.Sum256(data)
	lines := 0
	if len(data) > 0 {
		// Count a trailing unterminated line so summaries match editor-visible lines.
		lines = strings.Count(string(data), "\n")
		if data[len(data)-1] != '\n' {
			lines++
		}
	}
	return fmt.Sprintf("sha256:%s bytes:%d lines:%d", hex.EncodeToString(sum[:])[:16], len(data), lines)
}
