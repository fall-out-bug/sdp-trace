package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

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
