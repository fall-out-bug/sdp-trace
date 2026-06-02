package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func buildPRInputFromOptions(opts *flagSet) (packet.GitHubPREvidenceInput, error) {
	source, event, err := loadPRInputSourceEvent(opts)
	if err != nil {
		return packet.GitHubPREvidenceInput{}, err
	}
	// Event metadata seeds the input before optional local and live evidence is
	// layered in.
	input := githubPRInputFromEvent(event, source, os.Getenv)
	if err := completePRInputFromOptions(opts, source, &input); err != nil {
		// Optional evidence failures still invalidate the whole packet input,
		// because partial PR packets can overstate route or CI readiness.
		return packet.GitHubPREvidenceInput{}, err
	}
	return input, nil
}

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

func validPRInputSource(source string) bool {
	return source == "github-actions" || source == "github-fixture"
}

func prEventPath(source string, eventPath string, getenv func(string) string) string {
	if source == "github-actions" && eventPath == "" {
		// GitHub Actions events default to the runner-provided event file when
		// the packet command is not using an explicit fixture path.
		return getenv("GITHUB_EVENT_PATH")
	}
	return eventPath
}
