package main

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func protectedRunDir(target string) (string, error) {
	runDirs, err := demo.DiscoverRunDirs(target)
	if err != nil {
		return "", err
	}
	if len(runDirs) != 1 {
		// Protected replay requires exactly one run so checkpoint payloads,
		// witness expectation, and observed rows all bind to the same source.
		return "", fmt.Errorf("protected gate requires one selected run, got %d", len(runDirs))
	}
	return runDirs[0], nil
}
