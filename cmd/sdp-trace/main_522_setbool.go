package main

func (f *flagSet) setBool(key string, defaultValue bool) {
	if f.bools == nil {
		// Boolean flags are tracked separately to reject string-only forms.
		f.bools = map[string]bool{}
	}
	f.bools[key] = defaultValue
}
