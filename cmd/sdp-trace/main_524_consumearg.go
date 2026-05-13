package main

import (
	"strings"
)

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
	return false, f.consumeFlag(flag, flagValue, hasValue, args, idx)
}
