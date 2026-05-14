package main

import (
	"errors"
	"fmt"
	"os"
)

func loadPRInputSourceEvent(opts *flagSet) (string, prFixtureEvent, error) {
	source := opts.stringValue("source")
	if !validPRInputSource(source) {
		// Unknown sources are rejected instead of silently falling back to a live
		// runner mode that could cite the wrong event.
		return "", prFixtureEvent{}, fmt.Errorf("unsupported packet build-pr source %q", source)
	}
	// The source mode decides whether the event path is explicit fixture data or
	// the GitHub Actions event file.
	eventPath := prEventPath(source, opts.stringValue("github-event"), os.Getenv)
	if eventPath == "" {
		// Without an event file there is no authoritative PR identity for the
		// packet rows.
		return "", prFixtureEvent{}, errors.New("missing GitHub event JSON")
	}
	event, err := loadPRFixtureEvent(eventPath)
	if err != nil {
		// Bad fixture/event JSON is a source failure, not a partially verified
		// packet.
		return "", prFixtureEvent{}, err
	}
	return source, event, nil
}
