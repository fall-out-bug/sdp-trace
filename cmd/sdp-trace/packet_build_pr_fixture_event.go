package main

import (
	"errors"
	"strings"
)

type prFixtureEvent struct {
	WorkflowRunID string `json:"workflow_run_id"`
	PullRequest   struct {
		Number  int    `json:"number"`
		HTMLURL string `json:"html_url"`
		Title   string `json:"title"`
		BodyRef string `json:"body_ref"`
		DiffURL string `json:"diff_url"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
		Base struct {
			Ref string `json:"ref"`
			SHA string `json:"sha"`
		} `json:"base"`
		Head struct {
			Ref string `json:"ref"`
			SHA string `json:"sha"`
		} `json:"head"`
	} `json:"pull_request"`
}

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
