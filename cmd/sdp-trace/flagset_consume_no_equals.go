package main

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
