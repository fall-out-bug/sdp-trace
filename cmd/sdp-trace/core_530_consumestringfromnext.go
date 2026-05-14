package main

import (
	"fmt"
	"strings"
)

func (f *flagSet) consumeStringFromNext(flag string, args []string, idx *int) error {
	nextIdx := *idx + 1
	if nextIdx >= len(args) {
		// String flags must have a concrete following value.
		return fmt.Errorf("flag --%s requires value", flag)
	}
	value := args[nextIdx]
	if strings.HasPrefix(value, "--") {
		// Another flag cannot stand in for a missing string value.
		return fmt.Errorf("flag --%s requires value", flag)
	}
	*idx = nextIdx
	f.data[flag] = value
	return nil
}
