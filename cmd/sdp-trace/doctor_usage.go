package main

import (
	"fmt"
	"io"
)

const usageText = usagePrimaryText + usageTrustText + usagePacketText

func printUsage(w io.Writer) {
	// Global help is the canonical local command contract for this small CLI.
	fmt.Fprint(w, usageText)
}
