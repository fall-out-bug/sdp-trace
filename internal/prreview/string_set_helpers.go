package prreview

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
)

// String sets are sanitized before dedupe and sorting.
//
// This keeps unsafe variants from surviving in reason/action collections while
// producing deterministic output for validation and summary artifacts.
func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		out = appendUniqueSafeString(out, seen, value)
	}
	sort.Strings(out)
	return out
}

func appendUniqueSafeString(out []string, seen map[string]bool, value string) []string {
	value = safeText(value)
	if value == "" || seen[value] {
		return out
	}
	seen[value] = true
	return append(out, value)
}

func commandDigest(command []string) string {
	if len(command) == 0 {
		return ""
	}
	sum := sha256.Sum256([]byte(strings.Join(command, "\x00")))
	return "sha256:" + hex.EncodeToString(sum[:])
}
