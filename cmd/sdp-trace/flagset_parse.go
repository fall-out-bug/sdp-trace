package main

import (
	"fmt"
	"strings"
)

// parse processes args and populates the flagSet. Positional arguments are
// stored in f.args. Parsing stops at "--" and everything after it is treated
// as positional.
func (f *flagSet) parse(args []string) error {
	rest := make([]string, 0)
	for i := 0; i < len(args); i++ {
		// The loop index is passed by pointer so string flags can consume their
		// following value without reparsing it as positional input.
		// consumeArg owns index advancement for flags with following values.
		done, err := f.consumeArg(args, &i, &rest)
		if err != nil {
			return err
		}
		if done {
			f.args = rest
			return nil
		}
	}
	f.args = rest
	return nil
}

// consumeArg handles a single argument at *idx. It returns true when parsing
// should stop (because "--" was encountered), and false otherwise.
func (f *flagSet) consumeArg(args []string, idx *int, rest *[]string) (bool, error) {
	arg := args[*idx]
	if arg == "--" {
		// Everything after -- is command payload, not parser flags.
		*rest = append(*rest, args[*idx+1:]...)
		return true, nil
	}
	if !strings.HasPrefix(arg, "--") {
		// Positional arguments are preserved for command-specific validation.
		*rest = append(*rest, arg)
		return false, nil
	}
	flag, flagValue, hasValue := splitFlag(arg)
	// Flag value consumption may advance idx when the value is in the next arg.
	return false, f.consumeFlag(flag, flagValue, hasValue, args, idx)
}

// splitFlag extracts the flag name and optional value from an argument of the
// form "--flag" or "--flag=value".
func splitFlag(arg string) (string, string, bool) {
	parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
	if len(parts) == 1 {
		// Bare flags may either be boolean flags or string flags with next value.
		return parts[0], "", false
	}
	return parts[0], parts[1], true
}

// consumeFlag dispatches a flag to the appropriate value consumer based on its
// registration type and whether a value was provided inline.
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
