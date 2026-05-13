package main

import "io"

func errorOutput(opts options) io.Writer {
	if opts.err == nil {
		// Tests and embedded callers may omit stderr; threshold checks should
		// still return verdicts without forcing diagnostic output.
		return io.Discard
	}
	return opts.err
}
