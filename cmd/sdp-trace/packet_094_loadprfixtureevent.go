package main

import (
	"errors"
	"strings"
)

func loadPRFixtureEvent(path string) (prFixtureEvent, error) {
	var event prFixtureEvent
	if err := readOptionalJSON(path, &event); err != nil {
		return event, err
	}
	if event.PullRequest.Number == 0 || strings.TrimSpace(event.PullRequest.HTMLURL) == "" {
		// Fixture events need enough PR identity to build packet rows.
		return event, errors.New("missing pull_request metadata in GitHub event")
	}
	// Remaining PR fields may be empty in focused fixtures; row validation will
	// lower missing evidence as needed.
	return event, nil
}
