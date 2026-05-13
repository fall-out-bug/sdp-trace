package prreview

import (
	"errors"
	"fmt"
)

func validatePacketIdentityOptions(opts PacketOptions) error {
	// validatePacketIdentityOptions keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if !repoIDPattern.MatchString(opts.RepoID) {
		return fmt.Errorf("unsafe_repo_id: repo_id must match %s", repoIDPattern.String())
	}
	if !changeRefPattern.MatchString(opts.ChangeRef) {
		return fmt.Errorf("unsafe_change_ref: change_ref must match %s", changeRefPattern.String())
	}
	if !validPacketCommits(opts) {
		return errors.New("invalid_commit_sha: base and head must be 40 lowercase hex characters")
	}
	return nil
}
