package main

func (f *flagSet) isKnownFlag(flag string) (bool, bool) {
	// Return both flag classes so parsing can reject value syntax for booleans
	// without losing the unknown-flag distinction.
	_, isString := f.data[flag]
	_, isBool := f.bools[flag]
	return isString, isBool
}
