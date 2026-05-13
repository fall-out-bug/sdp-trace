package main

import (
	"errors"
	"strings"
)

func requirePRReviewPacketInputs(opts *flagSet) error {
	for _, flag := range prReviewPacketRequiredFlags {
		if strings.TrimSpace(opts.stringValue(flag.name)) == "" {
			// Required packet fields are provenance anchors, not cosmetic labels.
			return errors.New(flag.message)
		}
	}
	return nil
}
