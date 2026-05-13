package interaction

import (
	"strings"
	"time"
)

func normalizeImport(opts ImportOptions) ImportOptions {
	// normalizeImport keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	opts.TaskID = strings.TrimSpace(opts.TaskID)
	opts.Source = strings.TrimSpace(opts.Source)
	opts.SourceRef = strings.TrimSpace(opts.SourceRef)
	opts.EventsJSONL = strings.TrimSpace(opts.EventsJSONL)
	opts.Out = strings.TrimSpace(opts.Out)
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	return opts
}
