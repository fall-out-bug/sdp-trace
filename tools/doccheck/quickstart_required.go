package main

// requiredQuickstartCommands are the minimal commands that must appear in the
// contributor quick start so a new reader can verify the local environment.
// This slice is read-only; do not modify at runtime.
var requiredQuickstartCommands = []string{
	"go run ./cmd/sdp-trace --help",
	"go run ./cmd/sdp-trace doctor",
	"go run ./cmd/sdp-trace wrap",
	"go run ./cmd/sdp-trace verify",
	"go run ./cmd/sdp-trace explain",
}
