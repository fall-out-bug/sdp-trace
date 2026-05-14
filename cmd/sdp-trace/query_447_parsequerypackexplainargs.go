package main

import (
	"fmt"
	"strings"
)

func parseQueryPackExplainArgs(args []string) (*queryPackExplainOptions, error) {
	opts := &flagSet{name: "query-pack explain"}
	// Explanation takes one persisted result artifact; it has no run-directory
	// fallback.
	opts.setString("result", "")
	if err := opts.parse(args); err != nil {
		return nil, err
	}
	if len(opts.rest()) != 0 {
		return nil, fmt.Errorf("query-pack explain accepts only flags")
	}
	resultPath := strings.TrimSpace(opts.stringValue("result"))
	if resultPath == "" {
		// Explanation is bound to a persisted result artifact, not stdin or a
		// transient in-memory pack.
		return nil, fmt.Errorf("query-pack explain requires --result")
	}
	return &queryPackExplainOptions{resultPath: resultPath}, nil
}
