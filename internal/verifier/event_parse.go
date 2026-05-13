package verifier

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func appendParsedRunEvent(events []trace.Event, path string) ([]trace.Event, error) {
	event, err := parseRunEvent(path)
	if err != nil {
		// Surface only the basename to keep verifier errors portable.
		return nil, fmt.Errorf("invalid event %s: %w", filepath.Base(path), err)
	}
	return append(events, event), nil
}

func parseRunEvent(path string) (trace.Event, error) {
	var event trace.Event
	if err := trace.ReadJSON(context.Background(), path, &event); err != nil {
		return trace.Event{}, err
	}
	// Structural event validation happens during chain verification.
	return event, nil
}
