package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/releaseproof"
)

func evaluateAndWriteReleaseProof(opts *flagSet, stderr io.Writer) (releaseproof.Verification, int, bool) {
	repoRoot, err := releaseproof.RepoRoot(".")
	if err != nil {
		// Without a repository root, the manifest cannot be tied to an immutable
		// source boundary.
		fmt.Fprintln(stderr, err)
		return releaseproof.Verification{}, exitCannotVerify, false
	}
	return writeReleaseProofResult(repoRoot, opts, stderr)
}
