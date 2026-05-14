package main

func isHelpArg(arg string) bool {
	return arg == "help" || arg == "--help" || arg == "-h"
}
