package main

import (
	"fmt"
	"strings"
)

func (f *flagSet) consumeValue(flag, flagValue string, isBool bool) error {
	if !isBool {
		// --flag=value is the direct string assignment form.
		f.data[flag] = flagValue
		return nil
	}
	return f.consumeBoolValue(flag, flagValue)
}

func (f *flagSet) consumeNoEqualsValue(flag string, args []string, idx *int, isBool bool) error {
	if !isBool {
		// String flags without equals consume the next argument as their value.
		return f.consumeStringFromNext(flag, args, idx)
	}
	nextIdx := *idx + 1
	if !isBoolValueAt(args, nextIdx) {
		// Bare booleans default to true unless a literal follows the flag.
		f.bools[flag] = true
		return nil
	}
	*idx = nextIdx
	return f.consumeBoolValue(flag, args[*idx])
}

func (f *flagSet) consumeStringFromNext(flag string, args []string, idx *int) error {
	nextIdx := *idx + 1
	if nextIdx >= len(args) {
		// Missing string values are usage errors, not empty defaults.
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
