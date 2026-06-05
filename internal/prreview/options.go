package prreview

import "time"

type PacketOptions struct {
	OutDir            string
	RepoID            string
	ChangeRef         string
	BaseCommit        string
	HeadCommit        string
	DiffPath          string
	MetadataPath      string
	ContextPaths      []string
	VerificationPaths []string
	CIState           string
	CreatedBy         string
	Now               time.Time
}

// RunOptions contains local execution policy; it is intentionally separate
// from the portable packet so packet evidence stays harness-neutral.
type RunOptions struct {
	OutDir            string
	PacketDir         string
	AllowedRunners    map[string]bool
	Preview           bool
	Now               time.Time
	WorkDir           string
	NotAssessedReason string
}
