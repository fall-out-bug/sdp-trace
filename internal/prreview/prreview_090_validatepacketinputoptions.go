package prreview

import (
	"errors"
	"fmt"

	"strings"
)

func validatePacketInputOptions(opts PacketOptions) error {
	// validatePacketInputOptions keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if strings.TrimSpace(opts.DiffPath) == "" {

		return errors.New("pr_review_packet_requires_diff")
	}
	if invalidPacketCIState(opts) {

		return fmt.Errorf("invalid_ci_state: %s", opts.CIState)
	}
	return nil
}
