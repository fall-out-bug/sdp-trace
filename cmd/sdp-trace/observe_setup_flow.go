package main

import (
	"fmt"
	"io"
)

var runObserveSetupCommand = func(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseObserveSetupArgs(args, stderr)
	if !ok {
		return code
	}
	session, err := setupObservedSession(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeObserveSetup(stdout, stderr, session)
}
