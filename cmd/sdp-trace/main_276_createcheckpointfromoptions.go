package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func createCheckpointFromOptions(opts *flagSet, stderr io.Writer) (checkpoint.SignedCheckpoint, int, bool) {
	var key checkpoint.KeyPair
	// The private key file is local signing material, not an authority proof by
	// itself; policy binding is checked later during verification.
	if err := readJSONFile(opts.stringValue("private-key"), &key); err != nil {
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, 1, false
	}
	return createAndWriteCheckpoint(opts, key, stderr)
}
