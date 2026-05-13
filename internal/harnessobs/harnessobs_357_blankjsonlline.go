package harnessobs

import (
	"strings"
)

func blankJSONLLine(line []byte) bool {
	return len(strings.TrimSpace(string(line))) == 0
}
