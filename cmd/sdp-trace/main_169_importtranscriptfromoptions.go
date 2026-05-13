package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/interaction"
)

func importTranscriptFromOptions(opts *flagSet) (interaction.Trace, error) {
	// Source identity and source ref are preserved so imported transcript events
	// remain attributable after normalization.
	return interaction.ImportTranscript(interaction.ImportOptions{
		TaskID:      opts.stringValue("task-id"),
		Source:      opts.stringValue("source"),
		SourceRef:   opts.stringValue("source-ref"),
		EventsJSONL: opts.stringValue("events-jsonl"),
		Out:         opts.stringValue("out"),
	})
}
