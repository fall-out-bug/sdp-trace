package harnessobs

import (
	"fmt"

	"strings"
)

func rawScalarSignals(value any) []string {
	return []string{strings.ToLower(fmt.Sprint(value))}
}
