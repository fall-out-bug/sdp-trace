package prreview

import (
	"time"
)

func normalizeRunOptions(opts RunOptions) RunOptions {
	// normalizeRunOptions keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if opts.Now.IsZero() {

		opts.Now = time.Now().UTC()
	}
	if opts.WorkDir == "" {

		opts.WorkDir = "."
	}
	return opts
}
