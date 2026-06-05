package harnessobs

func extractCommandModelArgs(args []string) string {
	for i, arg := range args {
		if model, matched := commandModelArg(args, i, arg); matched {
			return model
		}
	}
	return ""
}
