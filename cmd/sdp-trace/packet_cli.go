package main

import (
	"fmt"
	"io"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func runPacketBuildGitHub(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePacketBuildGitHubOptions(args, stderr)
	if !ok {
		return code
	}
	// GitHub packet input is already materialized JSON; this command only
	// converts it into a portable packet bundle.
	input, err := packet.LoadGitHubInput(opts.stringValue("github-input"))
	if err != nil {
		// Input-read failure is reported as cannot_verify, not a packet fail.
		fmt.Fprintf(stderr, "read github input: %v\n", err)
		return exitCannotVerify
	}
	// Build and validate use a single clock tick for this local conversion.
	bundle := packet.BuildFromGitHubInput(input, time.Now().UTC())
	result := packet.Validate(bundle, time.Now().UTC())
	if result.State != packet.StatePass {
		// Invalid generated bundles are emitted as structured validation evidence.
		writeJSONPayloadUnchecked(stdout, result)
		return exitCannotVerify
	}
	// A passing validation result is the authority to persist the bundle.
	// Only validated bundles are written as durable packet artifacts.
	return writePacketBundle(opts.stringValue("out"), bundle, stdout, stderr)
}

func parsePacketBuildGitHubOptions(args []string, stderr io.Writer) (*flagSet, int, bool) {
	return parsePacketRequiredOptions(args, stderr, "packet build-github", "packet build-github accepts only flags", packetBuildGitHubRequiredFlags)
}

func writePacketBundle(outPath string, bundle packet.Bundle, stdout, stderr io.Writer) int {
	if err := writeJSONFile(outPath, bundle); err != nil {
		// Packet publication requires durable JSON, not stdout-only success.
		fmt.Fprintf(stderr, "write packet bundle: %v\n", err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "wrote %s\n", outPath)
	return 0
}

func runPacketValidate(args []string, stdout, stderr io.Writer) int {
	// Validate mode is intentionally narrow: one bundle in, one verdict out.
	opts := &flagSet{name: "packet validate"}
	opts.setString("bundle", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Bundle validation accepts only explicit artifact paths.
	if !requireOnlyFlags(opts, stderr, "packet validate accepts only flags", packetValidateRequiredFlags) {
		return exitUsage
	}
	// Validation reads a committed bundle and writes the verdict JSON unchanged.
	// It does not rebuild packet content from ambient repository state.
	bundle, err := packet.LoadBundle(opts.stringValue("bundle"))
	if err != nil {
		// An unreadable bundle means the packet verdict cannot be replayed.
		fmt.Fprintf(stderr, "read packet bundle: %v\n", err)
		return exitCannotVerify
	}
	// Validation time is local observation metadata for this CLI invocation.
	result := packet.Validate(bundle, time.Now().UTC())
	// Always publish the structured packet verdict before mapping the exit code.
	// Consumers should inspect the JSON state, not infer details from status.
	// The status code is only a shell-friendly summary of that verdict.
	writeJSONPayloadUnchecked(stdout, result)
	return packetValidationExit(result)
}

func runPacketCheckDemo(args []string, stdout, stderr io.Writer) int {
	// Demo-check mode applies the repository's first-packet readiness contract.
	opts := &flagSet{name: "packet check-demo"}
	opts.setString("bundle", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Demo checks are flag-only so no untracked positional input is trusted.
	if !requireOnlyFlags(opts, stderr, "packet check-demo accepts only flags", packetCheckDemoRequiredFlags) {
		return exitUsage
	}
	// The demo gate is a stricter first-packet contract over the same bundle.
	// It never consults live GitHub or CI state.
	bundle, err := packet.LoadBundle(opts.stringValue("bundle"))
	if err != nil {
		// Demo checks cannot infer trust from a missing packet bundle.
		fmt.Fprintf(stderr, "read packet bundle: %v\n", err)
		return exitCannotVerify
	}
	// Demo check time is local observation metadata for this CLI invocation.
	result := packet.CheckDemoFirstPacket(bundle, time.Now().UTC())
	// Demo gate output is intentionally the same validation envelope shape.
	// The exit code only distinguishes pass from expected demo-gate failure.
	// The JSON payload remains the detailed evidence for reviewers.
	writeJSONPayloadUnchecked(stdout, result)
	return packetDemoGateExit(result)
}

func runPacketRender(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePacketRenderOptions(args, stderr)
	if !ok {
		return code
	}
	// Rendering is read-only until the Markdown body has been generated.
	markdown, err := renderPacketBundleMarkdown(opts.stringValue("bundle"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// The generated markdown becomes evidence only after atomic persistence.
	return writePacketMarkdown(opts.stringValue("out"), markdown, stdout, stderr)
}

func renderPacketBundleMarkdown(bundlePath string) (string, error) {
	bundle, err := packet.LoadBundle(bundlePath)
	if err != nil {
		// Preserve the packet-read error as the trust boundary for render output.
		return "", fmt.Errorf("read packet bundle: %v", err)
	}
	markdown, err := packet.RenderMarkdown(bundle)
	if err != nil {
		// Markdown rendering may fail even when the bundle is readable.
		return "", fmt.Errorf("render packet: %v", err)
	}
	return markdown, nil
}

func parsePacketRenderOptions(args []string, stderr io.Writer) (*flagSet, int, bool) {
	return parsePacketRequiredOptions(args, stderr, "packet render", "packet render accepts only flags", packetRenderRequiredFlags)
}

func parsePacketRequiredOptions(args []string, stderr io.Writer, name, restMessage string, required []requiredCLIFlag) (*flagSet, int, bool) {
	opts := &flagSet{name: name}
	for _, flag := range required {
		// Packet commands share the same required string-flag parser.
		opts.setString(flag.name, "")
	}
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Positional payload is not accepted by packet helper commands.
	if !requireOnlyFlags(opts, stderr, restMessage, required) {
		return nil, exitUsage, false
	}
	// A successful parse means every required artifact path is non-empty.
	return opts, 0, true
}

func writePacketMarkdown(outPath, markdown string, stdout, stderr io.Writer) int {
	if err := writeTextFileAtomic(outPath, markdown); err != nil {
		// Markdown packets are written atomically to avoid partial review docs.
		fmt.Fprintf(stderr, "write packet: %v\n", err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "wrote %s\n", outPath)
	return 0
}
