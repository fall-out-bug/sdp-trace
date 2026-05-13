package main

import (
	"errors"
	"strings"
)

func requirePRReviewCheckInputs(opts *flagSet) error {
	outDir := opts.stringValue("out")
	if strings.TrimSpace(outDir) == "" {
		// A combined review check needs a directory because it writes multiple
		// artifacts whose paths become later evidence refs.
		return errors.New("pr-review check requires --out")
	}
	return requirePRReviewPacketInputs(opts)
}
