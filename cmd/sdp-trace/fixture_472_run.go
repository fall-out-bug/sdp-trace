package main

import (
	"context"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func runValidateFixtures(_ context.Context, args []string, stdout, stderr io.Writer) int {
	fixtureRoot := fixtureRootArg(args)
	// Fixture discovery is rooted explicitly so validation cannot wander into
	// unrelated run artifacts.
	runDirs, err := demo.DiscoverRunDirs(fixtureRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if validateFixtureRuns(fixtureRoot, runDirs, stdout, stderr) {
		return 1
	}
	return 0
}
