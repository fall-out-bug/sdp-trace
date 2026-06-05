package harnessobs

func commandModelArg(args []string, i int, arg string) (string, bool) {
	if arg == "--model" || arg == "-m" {
		return nextCommandModelArg(args, i), true
	}
	if model, ok := prefixedCommandModelArg(arg); ok {
		return model, true
	}
	return "", false
}
