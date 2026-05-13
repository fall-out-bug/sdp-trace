package main

import (
	"io"
)

func validateFixtureRuns(fixtureRoot string, runDirs []string, stdout, stderr io.Writer) bool {
	failed := false
	for _, runDir := range runDirs {
		// Continue through all fixtures so one broken run does not hide other
		// drift in the example corpus.
		if validateFixtureRun(fixtureRoot, runDir, stdout, stderr) {
			failed = true
		}
	}
	return failed
}
