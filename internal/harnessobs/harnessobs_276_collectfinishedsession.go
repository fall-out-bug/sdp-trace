package harnessobs

import (
	"time"
)

func collectFinishedSession(opts SessionOptions, session SessionRun, waitErr error, end time.Time) (SessionRun, Run, error) {
	// collectFinishedSession keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := writeFinishedSession(opts.OutDir, &session, end); err != nil {
		return SessionRun{}, Run{}, err
	}

	collected, observed, err := CollectSession(SessionCollectOptions{ProfilePath: opts.ProfilePath, RunDir: opts.OutDir, Now: end})
	if err != nil {
		return SessionRun{}, Run{}, err
	}
	if waitErr != nil {
		return collected, observed, waitErr
	}
	return collected, observed, nil
}
