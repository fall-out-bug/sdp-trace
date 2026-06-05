package prreview

import (
	"errors"
	"fmt"
	"strings"
)

// validatePacketOptions keeps packet generation inputs bounded before any files
// are written into the reviewer packet.
func validatePacketOptions(opts PacketOptions) error {
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

// validatePacketInputOptions checks local reviewer inputs that are not part of
// the immutable repository/change identity.
func validatePacketInputOptions(opts PacketOptions) error {
	if strings.TrimSpace(opts.DiffPath) == "" {
		return errors.New("pr_review_packet_requires_diff")
	}
	if invalidPacketCIState(opts) {
		return fmt.Errorf("invalid_ci_state: %s", opts.CIState)
	}
	return nil
}

// invalidPacketCIState treats an omitted CI state as not_assessed while still
// rejecting unsupported explicit state names.
func invalidPacketCIState(opts PacketOptions) bool {
	return opts.CIState != "" && !validCIState(opts.CIState)
}
