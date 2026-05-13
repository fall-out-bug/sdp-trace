package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func runReport(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts, target, code, ok := parseReportArgs(args, stderr)
	if !ok {
		return code
	}
	// Reports regenerate summary artifacts from run evidence; stdout mirrors
	// the summary but the output directory is the durable review surface.
	artifacts, err := demo.WriteReport(target, opts.stringValue("out"), opts.stringValue("contract"))
	if err != nil {
		// Report generation failure means no durable summary can be trusted.
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(artifacts.Summary, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return reportExitCode(artifacts.Summary)
}
