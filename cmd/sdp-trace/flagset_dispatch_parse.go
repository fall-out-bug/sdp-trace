package main

import (
	"fmt"
	"strings"
)

func splitFlag(arg string) (string, string, bool) {
	parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
	if len(parts) == 1 {
		// Bare flags may be boolean flags or string flags with next value.
		return parts[0], "", false
	}
	return parts[0], parts[1], true
}

func (f *flagSet) consumeFlag(flag string, flagValue string, hasValue bool, args []string, idx *int) error {
	isString, isBool := f.isKnownFlag(flag)
	if !isString && !isBool {
		// Unknown flags fail early before command code interprets inputs.
		return fmt.Errorf("unknown flag --%s", flag)
	}
	if hasValue {
		return f.consumeValue(flag, flagValue, isBool)
	}
	return f.consumeNoEqualsValue(flag, args, idx, isBool)
}
