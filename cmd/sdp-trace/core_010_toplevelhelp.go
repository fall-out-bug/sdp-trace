package main

func topLevelHelp(args []string) bool {
	return len(args) == 0 || isHelpArg(args[0])
}
