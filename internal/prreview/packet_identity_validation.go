package prreview

import (
	"errors"
	"fmt"
)

// validatePacketIdentityOptions checks portable source identity before packet
// construction records it as review context.
func validatePacketIdentityOptions(opts PacketOptions) error {
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

// validPacketCommits requires immutable 40-character lowercase SHA-1 refs for
// both sides of the reviewed change.
func validPacketCommits(opts PacketOptions) bool {
	return sha40Pattern.MatchString(opts.BaseCommit) && sha40Pattern.MatchString(opts.HeadCommit)
}
