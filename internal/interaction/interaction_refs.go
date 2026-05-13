package interaction

import (
	"path/filepath"
	"strings"
)

func validReference(ref string) bool {
	// validReference keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if strings.TrimSpace(ref) == "" || strings.Contains(ref, " ") {
		return false
	}
	if strings.HasPrefix(ref, "recorder:") {
		return recorderRefPattern.MatchString(ref)
	}
	return contentRefPattern.MatchString(ref)
}

func contentBlobPath(tracePath string, event Event) string {

	return filepath.Join(filepath.Dir(tracePath), "interactions", event.TaskID, event.InteractionID+".txt")
}
