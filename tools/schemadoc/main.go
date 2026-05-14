package main

// main is the entry point for the schemadoc tool.
// It supports three modes:
//   - default: check schema/index.json against the schema directory
//   - -generate: print the Markdown table to stdout
//   - -verify-readme: ensure README.md table matches the index

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var g, v bool
	flag.BoolVar(&g, "generate", false, "print generated README table section to stdout")
	flag.BoolVar(&v, "verify-readme", false, "verify README.md schema table matches index.json")
	flag.Parse()
	os.Exit(exitCode(run(g, v), os.Stderr))
}

func exitCode(err error, stderr io.Writer) int {
	if err == nil {
		return 0
	}
	fmt.Fprintln(stderr, err)
	return 1
}
