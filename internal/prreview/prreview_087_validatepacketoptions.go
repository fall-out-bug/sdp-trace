package prreview

import (
	"errors"

	"strings"
)

func validatePacketOptions(opts PacketOptions) error {
	// validatePacketOptions keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if strings.TrimSpace(opts.OutDir) == "" {

		return errors.New("pr_review_packet_requires_out")
	}
	if err := validatePacketIdentityOptions(opts); err != nil {
		return err
	}
	if err := validatePacketInputOptions(opts); err != nil {
		return err
	}
	return nil
}
