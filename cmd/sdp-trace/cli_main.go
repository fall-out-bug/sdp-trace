package main

import (
	"io"
	"os"
)

var cliStdin io.Reader = os.Stdin

var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if topLevelHelp(args) {
		printUsage(stdout)
		return 0
	}
	// The first token is the only command selector; everything after it stays
	// command-owned so subcommands can preserve their own evidence contract.
	return dispatchCommand(args[0], args[1:], stdout, stderr)
}

func topLevelHelp(args []string) bool {
	return len(args) == 0 || isHelpArg(args[0])
}
