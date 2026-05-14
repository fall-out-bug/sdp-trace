package main

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/prreview"
)

func preparePRReviewCheck(opts *flagSet, args []string, stderr io.Writer) (prreview.Packet, prreview.ReviewProfile, int, bool) {
	outDir := opts.stringValue("out")
	// Packet construction is first because every later review artifact is bound
	// to the packet identity and change metadata.
	packet, err := prreview.BuildPacket(prReviewPacketOptions(opts, args, filepath.Join(outDir, "packet")))
	if err != nil {
		// Packet construction failure means no downstream review artifact should
		// be produced.
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, exitCannotVerify, false
	}
	// The profile controls review planes and runner policy; missing profile
	// evidence prevents the combined check from claiming review coverage.
	profile, err := prreview.ReadProfile(opts.stringValue("profile"))
	if err != nil {
		// A malformed profile means required planes and runners cannot be known.
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, exitCannotVerify, false
	}
	if err := requireDirectory(opts.stringValue("work-dir")); err != nil {
		// Review runners operate relative to a concrete working directory.
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, exitCannotVerify, false
	}
	// Profile and directory checks are part of preparation, not a review verdict.
	// A prepared check has all inputs needed for the runner boundary and no
	// persisted artifact has been written yet.
	return packet, profile, 0, true
}
