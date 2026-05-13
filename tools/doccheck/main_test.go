package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRunAcceptsCurrentCommandSurface(t *testing.T) {
	if err := run(); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestExitCode(t *testing.T) {
	var stderr bytes.Buffer
	if got := exitCode(nil, &stderr); got != 0 {
		t.Fatalf("exitCode(nil) = %d, want 0", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("exitCode(nil) wrote stderr: %q", stderr.String())
	}

	err := errors.New("drift")
	if got := exitCode(err, &stderr); got != 1 {
		t.Fatalf("exitCode(err) = %d, want 1", got)
	}
	if !strings.Contains(stderr.String(), "drift") {
		t.Fatalf("exitCode(err) stderr = %q, want drift", stderr.String())
	}
}

func TestCompareCommandSurfacePassesWhenDocMatchesHelp(t *testing.T) {
	help := strings.Join([]string{
		"sdp-trace local recorder and verifier commands.",
		"  sdp-trace version",
		"  sdp-trace verify <run-dir>",
	}, "\n")
	doc := strings.Join([]string{
		"Current command surface:",
		"",
		"- `sdp-trace --help`",
		"- `sdp-trace verify <run-dir>`",
		"- `sdp-trace version`",
		"",
		"Do not add aliases",
	}, "\n")
	if err := compareCommandSurface(help, doc); err != nil {
		t.Fatalf("compareCommandSurface: %v", err)
	}
}

func TestCompareCommandSurfaceReportsMissingAndStaleCommands(t *testing.T) {
	help := strings.Join([]string{
		"sdp-trace local recorder and verifier commands.",
		"  sdp-trace version",
		"  sdp-trace verify <run-dir>",
	}, "\n")
	doc := strings.Join([]string{
		"Current command surface:",
		"",
		"- `sdp-trace --help`",
		"- `sdp-trace verify <run-dir>`",
		"- `sdp-trace old-command`",
		"",
		"Do not add aliases",
	}, "\n")
	err := compareCommandSurface(help, doc)
	if err == nil {
		t.Fatal("compareCommandSurface succeeded, want drift error")
	}
	for _, want := range []string{"missing documented commands: sdp-trace version", "stale documented commands: sdp-trace old-command"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error missing %q: %v", want, err)
		}
	}
}
