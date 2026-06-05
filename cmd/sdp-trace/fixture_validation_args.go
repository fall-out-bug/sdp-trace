package main

func fixtureRootArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	// Default fixtures to the current directory for local demo validation.
	return "."
}
