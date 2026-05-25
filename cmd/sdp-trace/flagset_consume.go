package main

import (
	"fmt"
	"strings"
)

// consumeValue handles a flag that has an inline value ("--flag=value").
func (f *flagSet) consumeValue(flag, flagValue string, isBool bool) error {
	if !isBool {
		// --flag=value is the direct string assignment form.
		f.data[flag] = flagValue
		return nil
	}
	return f.consumeBoolValue(flag, flagValue)
}

// consumeNoEqualsValue handles a flag without an inline value ("--flag").
// String flags consume the next argument; boolean flags default to true unless
// the next argument is a boolean literal.
func (f *flagSet) consumeNoEqualsValue(flag string, args []string, idx *int, isBool bool) error {
	if !isBool {
		// String flags without equals consume the next argument as their value.
		return f.consumeStringFromNext(flag, args, idx)
	}
	nextIdx := *idx + 1
	if !isBoolValueAt(args, nextIdx) {
		// Bare boolean flags imply true unless followed by a boolean literal.
		f.bools[flag] = true
		return nil
	}
	*idx = nextIdx
	return f.consumeBoolValue(flag, args[*idx])
}

// consumeStringFromNext consumes the next positional argument as the value for
// a string flag.
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

// isBoolValueAt reports whether args[idx] is a boolean literal.
func isBoolValueAt(args []string, idx int) bool {
	return idx < len(args) && isBoolLiteral(args[idx])
}

// consumeBoolValue parses a boolean literal and stores it in the flagSet.
func (f *flagSet) consumeBoolValue(flag, flagValue string) error {
	switch strings.ToLower(flagValue) {
	case "false", "0":
		// Accept compact false spellings for generated command lines.
		f.bools[flag] = false
	case "true", "1", "":
		// Empty value covers --flag= and keeps legacy true semantics.
		f.bools[flag] = true
	default:
		// Invalid boolean values are usage errors, not ignored arguments.
		return fmt.Errorf("invalid boolean value for --%s: %s", flag, flagValue)
	}
	return nil
}
