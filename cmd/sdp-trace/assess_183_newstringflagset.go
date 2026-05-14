package main

func newStringFlagSet(name string, flags []string) *flagSet {
	opts := &flagSet{name: name}
	// Shared assess flags keep the command surface stable while each selected
	// profile validates only the inputs it can actually use.
	for _, flag := range flags {
		opts.setString(flag, "")
	}
	return opts
}
