package repoobserver

import "time"

func withDefaultNow(opts Options) Options {
	if opts.Now.IsZero() {
		// Repo observations use wall-clock UTC in normal runs, while tests inject
		// time for deterministic status output.
		opts.Now = time.Now().UTC()
	}
	return opts
}
