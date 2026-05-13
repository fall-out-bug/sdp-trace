package main

import (
	"fmt"
	"io"
)

// Packet subcommands share one flag-only parser so each artifact input is named
// before any packet bundle is built, validated, or rendered.

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
