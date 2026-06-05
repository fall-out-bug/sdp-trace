package main

import "strings"

func (f *flagSet) parse(args []string) error {
	rest := make([]string, 0)
	for i := 0; i < len(args); i++ {
		// Each iteration either stores positional input or delegates a complete
		// flag token, preserving argv order for command-specific validation.
		done, err := f.consumeArg(args, &i, &rest)
		if err != nil {
			return err
		}
		if done {
			// A stop signal means "--" moved the remaining payload into rest.
			f.args = rest
			return nil
		}
	}
	f.args = rest
	return nil
}

func (f *flagSet) consumeArg(args []string, idx *int, rest *[]string) (bool, error) {
	arg := args[*idx]
	if arg == "--" {
		// Everything after -- is command payload, not parser flags; callers use
		// the returned stop signal to avoid reparsing payload words.
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
