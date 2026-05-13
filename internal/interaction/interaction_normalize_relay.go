package interaction

import (
	"strings"
	"time"
)

func normalizeRelay(opts RelayOptions) RelayOptions {
	// normalizeRelay keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	opts.TaskID = strings.TrimSpace(opts.TaskID)
	opts.ActorType = strings.TrimSpace(opts.ActorType)
	opts.ActorID = strings.TrimSpace(opts.ActorID)
	opts.Target = strings.TrimSpace(opts.Target)
	opts.EventType = strings.TrimSpace(opts.EventType)
	opts.OperationID = strings.TrimSpace(opts.OperationID)
	opts.StageID = strings.TrimSpace(opts.StageID)
	opts.Out = strings.TrimSpace(opts.Out)
	return defaultRelayOptions(opts)
}

func defaultRelayOptions(opts RelayOptions) RelayOptions {
	// defaultRelayOptions keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	opts = defaultRelayStrings(opts)
	if opts.Now.IsZero() {
		opts.Now = time.Now().UTC()
	}
	return opts
}

func defaultRelayStrings(opts RelayOptions) RelayOptions {
	// defaultRelayStrings keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if opts.ActorType == "" {
		opts.ActorType = "human_user"
	}
	if opts.Target == "" {
		opts.Target = "agent"
	}
	if opts.EventType == "" {
		opts.EventType = "corrective_feedback"
	}
	return opts
}
