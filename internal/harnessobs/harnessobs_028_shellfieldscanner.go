package harnessobs

import (
	"strings"
)

type shellFieldScanner struct {
	fields  []string
	current strings.Builder
	quote   rune
	escaped bool
}
