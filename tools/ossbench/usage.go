package main

import (
	"flag"
	"fmt"
	"io"
)

func usageFunc(fs *flag.FlagSet, stderr io.Writer) func() {
	return func() {
		fmt.Fprintf(stderr, "Usage: ossbench [flags] [command args...]\n")
		fmt.Fprintf(stderr, "Run built-in or custom benchmarks with min/max/median stats.\n")
		fmt.Fprintf(stderr, "Use -list to see built-in benchmarks.\n\n")
		fs.PrintDefaults()
	}
}
