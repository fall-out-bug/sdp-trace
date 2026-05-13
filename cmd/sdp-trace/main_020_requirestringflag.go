package main

import (
	"fmt"
	"io"
	"strings"
)

func requireStringFlag(opts *flagSet, stderr io.Writer, flag, message string) bool {
	if strings.TrimSpace(opts.stringValue(flag)) != "" {
		return true
	}
	// Empty string flags are missing evidence inputs even if the flag appeared.
	fmt.Fprintln(stderr, message)
	return false
}
