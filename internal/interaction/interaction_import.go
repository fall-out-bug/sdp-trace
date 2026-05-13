package interaction

import (
	"errors"
	"strings"
)

func ImportTranscript(opts ImportOptions) (Trace, error) {
	// ImportTranscript keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	opts = normalizeImport(opts)
	if err := validateImportOptions(opts); err != nil {
		return Trace{}, err
	}
	events, err := importTranscriptEvents(opts)
	if err != nil {
		return Trace{}, err
	}
	trace := NewTrace(opts.TaskID, SourcePreclassifiedTranscript, events, opts.Now)
	if err := WriteTrace(opts.Out, trace); err != nil {
		return Trace{}, err
	}
	return trace, nil
}

func validateImportOptions(opts ImportOptions) error {
	// validateImportOptions keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if opts.Source != SourcePreclassifiedTranscript {
		return errors.New("interaction import-transcript requires --source preclassified-transcript-import")
	}
	if err := validateSafeID("task_id", opts.TaskID); err != nil {
		return err
	}
	if strings.TrimSpace(opts.EventsJSONL) == "" {
		return errors.New("interaction import-transcript requires --events-jsonl")
	}
	return nil
}

func importTranscriptEvents(opts ImportOptions) ([]Event, error) {
	// importTranscriptEvents keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	events, err := readJSONLEvents(opts.EventsJSONL)
	if err != nil {
		return nil, err
	}
	if err := normalizeTranscriptEvents(events, opts); err != nil {
		return nil, err
	}
	return events, nil
}
