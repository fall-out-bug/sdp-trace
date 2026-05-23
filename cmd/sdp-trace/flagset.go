package main

// flagSet is a minimal command-line flag parser used by subcommands that need
// lightweight string and boolean flag handling without importing the full
// standard flag package.
type flagSet struct {
	name  string
	data  map[string]string
	bools map[string]bool
	args  []string
}

// setString registers a string flag with its default value.
func (f *flagSet) setString(key, defaultValue string) {
	if f.data == nil {
		// Lazily allocate so tiny commands only register the flags they own.
		f.data = map[string]string{}
	}
	f.data[key] = defaultValue
}

// setBool registers a boolean flag with its default value.
func (f *flagSet) setBool(key string, defaultValue bool) {
	if f.bools == nil {
		// Boolean flags are tracked separately to reject string-only forms.
		f.bools = map[string]bool{}
	}
	f.bools[key] = defaultValue
}

// stringValue returns the parsed value for a registered string flag.
func (f *flagSet) stringValue(key string) string {
	if f.data == nil {
		// Unregistered string maps read as absent flags, matching parse defaults.
		return ""
	}
	return f.data[key]
}

// boolValue returns the parsed value for a registered boolean flag.
func (f *flagSet) boolValue(key string) bool {
	if f.bools == nil {
		// Unregistered bool maps default to false instead of implying a flag.
		return false
	}
	return f.bools[key]
}

// rest returns the positional arguments remaining after flag parsing.
func (f *flagSet) rest() []string {
	return f.args
}

// isKnownFlag reports whether flag is a registered string or boolean flag.
func (f *flagSet) isKnownFlag(flag string) (bool, bool) {
	// Return both flag classes so parsing can reject value syntax for booleans
	// without losing the unknown-flag distinction.
	_, isString := f.data[flag]
	_, isBool := f.bools[flag]
	return isString, isBool
}
