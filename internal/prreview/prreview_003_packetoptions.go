package prreview

import (
	"time"
)

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

// Packet is the frozen review input boundary: reviewers assess the copied refs,
// not the mutable files those refs came from.
