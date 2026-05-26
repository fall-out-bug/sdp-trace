package main

import (
	"fmt"
	"time"
)

func runIteration(def benchmarkDef, index int) (time.Duration, string, bool) {
	dur, err, timedOut := runSingleCommand(def.Cmd, def.Args, def.Dir)
	if err != nil {
		return 0, fmt.Sprintf("iteration %d failed: %v", index, err), timedOut
	}
	return dur, "", false
}
