package recorder

import (
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func (w *runWriter) persistEvent(event trace.Event, eventType trace.EventType) error {
	// Sequence-prefixed filenames keep the append order visible in filesystem
	// listings without making the filename part of the event hash.
	filename := filepath.Join(w.runDir, "events", eventFilename(event.Sequence, eventType))
	return writeIndentedJSON(filename, event)
}

func (w *runWriter) advanceEventHead(event trace.Event) {
	// The in-memory head mirrors the persisted event chain and feeds both the
	// next event's previous hash and the manifest closure fields.
	w.sequence++
	w.lastHash = event.EventHash
	w.events = append(w.events, event)
}
