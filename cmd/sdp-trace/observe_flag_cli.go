package main

import (
	"fmt"
	"io"
)

type observeStringFlag struct {
	name         string
	defaultValue string
}

func parseFlagOnlyCommand(args []string, stderr io.Writer, name, restMessage string, flags []observeStringFlag, required []requiredCLIFlag) (*flagSet, int, bool) {
	opts := &flagSet{name: name}
	for _, flag := range flags {
		opts.setString(flag.name, flag.defaultValue)
	}
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireOnlyFlags(opts, stderr, restMessage, required) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func writeJSONExit(stdout, stderr io.Writer, payload any, context string) int {
	if !writeJSONPayload(stdout, stderr, payload, context) {
		return exitCannotVerify
	}
	return 0
}

func runJSONFlagCommand[T any](args []string, stdout, stderr io.Writer, parse func([]string, io.Writer) (*flagSet, int, bool), build func(*flagSet) (T, error), context string) int {
	opts, code, ok := parse(args, stderr)
	if !ok {
		return code
	}
	payload, err := build(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeJSONExit(stdout, stderr, payload, context)
}
