package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func executePRReviewRun(opts *flagSet, args []string) (prreview.RunSet, *prreview.RunPreview, error) {
	packet, profile, err := readPRReviewPacketAndProfileValues(opts)
	if err != nil {
		return prreview.RunSet{}, nil, err
	}
	// The working directory is part of the runner boundary; nonexistent paths
	// would make external review evidence impossible to reproduce.
	if err := requireDirectory(opts.stringValue("work-dir")); err != nil {
		return prreview.RunSet{}, nil, err
	}
	// Runner allow-list values are reconstructed from raw args so repeated flags
	// cannot be collapsed by the flag parser.
	return prreview.RunReview(packet, profile, prreview.RunOptions{
		OutDir:            opts.stringValue("out"),
		AllowedRunners:    allowedRunnerSet(repeatedFlagValues(args, "allow-external-runner", opts.stringValue("allow-external-runner"))),
		Preview:           opts.boolValue("preview"),
		WorkDir:           opts.stringValue("work-dir"),
		NotAssessedReason: opts.stringValue("not-assessed-reason"),
	})
}
