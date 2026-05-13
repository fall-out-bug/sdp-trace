package harnessobs

import (
	"strings"
)

func safeTokenRune(r rune) bool {
	return strings.ContainsRune(safeTokenRunes, r)
}
