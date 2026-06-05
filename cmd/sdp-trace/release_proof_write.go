package main

import (
	"fmt"
	"io"
	"time"

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

func writeReleaseProofResult(repoRoot string, opts *flagSet, stderr io.Writer) (releaseproof.Verification, int, bool) {
	result, err := releaseproof.Evaluate(repoRoot, opts.stringValue("manifest"), time.Now())
	if err != nil {
		// Evaluation errors leave release proof unverifiable; they are not a
		// successful proof with warnings.
		fmt.Fprintln(stderr, err)
		return releaseproof.Verification{}, exitCannotVerify, false
	}
	if err := releaseproof.Write(opts.stringValue("out"), result); err != nil {
		// A proof that cannot be persisted cannot be referenced by later gates.
		fmt.Fprintln(stderr, err)
		return releaseproof.Verification{}, 1, false
	}
	return result, 0, true
}
