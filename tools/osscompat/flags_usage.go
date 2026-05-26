package main

import (
	"flag"
	"fmt"
	"io"
)

func usageFunc(fs *flag.FlagSet, stderr io.Writer) func() {
	return func() {
		fmt.Fprintf(stderr, "Usage: osscompat [flags]\n")
		fmt.Fprintf(stderr, "Run compatibility probes and emit results.\n")
		fmt.Fprintf(stderr, "Exit 0 means no probe returned fail; it does NOT mean all probes passed.\n\n")
		fs.PrintDefaults()
	}
}
