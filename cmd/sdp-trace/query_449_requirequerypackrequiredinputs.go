package main

import (
	"fmt"
)

func requireQueryPackRequiredInputs(runPath, outPath string) error {
	if runPath == "" {
		// The run path is the replayable source evidence for this pack.
		return fmt.Errorf("query-pack requires --run")
	}
	if outPath == "" {
		// The output path is the durable artifact reviewed by later commands.
		return fmt.Errorf("query-pack requires --out")
	}
	return nil
}
