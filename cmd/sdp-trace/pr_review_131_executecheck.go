package main

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func executePRReviewCheck(packet prreview.Packet, profile prreview.ReviewProfile, opts *flagSet, args []string, stderr io.Writer) (prreview.RunSet, *prreview.RunPreview, int, bool) {
	outDir := opts.stringValue("out")
	// Combined execution writes runs under a stable subdirectory so ledger and
	// validation artifacts can refer to the same run set.
	runs, preview, err := prreview.RunReview(packet, profile, prreview.RunOptions{
		OutDir:            filepath.Join(outDir, "runs"),
		PacketDir:         filepath.Join(outDir, "packet"),
		AllowedRunners:    allowedRunnerSet(repeatedFlagValues(args, "allow-external-runner", opts.stringValue("allow-external-runner"))),
		Preview:           opts.boolValue("preview"),
		WorkDir:           opts.stringValue("work-dir"),
		NotAssessedReason: opts.stringValue("not-assessed-reason"),
	})
	if err != nil {
		// Runner failures are recorded as cannot_verify because review evidence is
		// incomplete.
		fmt.Fprintln(stderr, err)
		return prreview.RunSet{}, nil, exitCannotVerify, false
	}
	return runs, preview, 0, true
}
