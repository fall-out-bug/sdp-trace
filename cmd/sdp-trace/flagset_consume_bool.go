package main

import (
	"fmt"
	"strings"
)

func (f *flagSet) consumeBoolValue(flag, flagValue string) error {
	switch strings.ToLower(flagValue) {
	case "false", "0":
		f.bools[flag] = false
	case "true", "1", "":
		// Empty value covers --flag= and keeps legacy true semantics.
		f.bools[flag] = true
	default:
		return fmt.Errorf("invalid boolean value for --%s: %s", flag, flagValue)
	}
	return nil
}
