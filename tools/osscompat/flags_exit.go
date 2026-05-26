package main

import "flag"

func flagParseExitCode(err error) int {
	if err == flag.ErrHelp {
		return 0
	}
	return 2
}
