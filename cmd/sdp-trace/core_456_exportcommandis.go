package main

func exportCommandIs(args []string, command string) bool {
	return len(args) > 0 && args[0] == command
}
