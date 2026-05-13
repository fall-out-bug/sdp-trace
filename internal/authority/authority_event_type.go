package authority

import (
	"strings"
)

func validEventType(event string) bool {
	return standardEventTypes[event] || strings.HasPrefix(event, "custom:")
}
