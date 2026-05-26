package main

import "fmt"

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
