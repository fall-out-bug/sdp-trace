package harnessobs

import (
	"strings"
)

func blankJSON(data []byte) bool {
	return strings.TrimSpace(string(data)) == ""
}
