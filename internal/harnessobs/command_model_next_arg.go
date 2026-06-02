package harnessobs

func nextCommandModelArg(args []string, i int) string {
	if i+1 >= len(args) {
		return ""
	}

	return safeCommandModel(args[i+1])
}
