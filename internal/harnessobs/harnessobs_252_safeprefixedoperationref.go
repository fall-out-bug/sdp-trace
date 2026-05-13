package harnessobs

import (
	"strings"
)

func safePrefixedOperationRef(ref string) bool {
	return !strings.Contains(ref, "..") && !strings.Contains(ref, "://") && len(ref) <= 256
}
