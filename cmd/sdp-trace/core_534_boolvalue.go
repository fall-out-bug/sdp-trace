package main

func (f *flagSet) boolValue(key string) bool {
	if f.bools == nil {
		// Unregistered bool maps default to false instead of implying a flag.
		return false
	}
	return f.bools[key]
}
