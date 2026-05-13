package main

func exportSubcommandIs(args []string, subcommand string) bool {
	return len(args) > 1 && args[1] == subcommand
}
