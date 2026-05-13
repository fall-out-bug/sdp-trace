package prreview

import (
	"strings"
)

func emptyReviewerOutput(output []byte) bool {
	return len(strings.TrimSpace(string(output))) == 0
}
