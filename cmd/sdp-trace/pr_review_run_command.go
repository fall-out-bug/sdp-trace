package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func runPRReviewRun(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePRReviewRunArgs(args, stderr)
	if !ok {
		return code
	}
	// Reviewer execution can only produce usable evidence when packet, profile,
	// runner allow-list, and work directory are all replayable.
	runs, preview, err := executePRReviewRun(opts, args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	writePRReviewRunOutput(stdout, runs, preview)
	return 0
}

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
		PacketDir:         packetDir(opts.stringValue("packet")),
		AllowedRunners:    allowedRunnerSet(repeatedFlagValues(args, "allow-external-runner", opts.stringValue("allow-external-runner"))),
		Preview:           opts.boolValue("preview"),
		WorkDir:           opts.stringValue("work-dir"),
		NotAssessedReason: opts.stringValue("not-assessed-reason"),
	})
}
