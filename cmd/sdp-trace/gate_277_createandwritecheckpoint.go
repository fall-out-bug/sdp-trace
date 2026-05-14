package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func createAndWriteCheckpoint(opts *flagSet, key checkpoint.KeyPair, stderr io.Writer) (checkpoint.SignedCheckpoint, int, bool) {
	// checkpoint.Create binds the requested run directory and signer identity
	// into the signed payload before this CLI writes the JSON artifact.
	// The CLI does not infer signer authority here; policy binding is replayed
	// by checkpoint verification.
	created, err := checkpoint.Create(opts.stringValue("run"), checkpoint.CreateOptions{
		CheckpointID: opts.stringValue("id"),
		SignerID:     opts.stringValue("signer-id"),
		Key:          key,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, 1, false
	}
	if err := writeJSONFile(opts.stringValue("out"), created); err != nil {
		// A checkpoint that cannot be persisted is not usable evidence.
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, 1, false
	}
	// Return the created artifact only after it exists at the requested path.
	return created, 0, true
}
