package harnessobs

import (
	"strings"
)

func operationRefPrefix(ref string) bool {
	return strings.HasPrefix(ref, "adapter-run:") || strings.HasPrefix(ref, "delivery-trace:")
}
