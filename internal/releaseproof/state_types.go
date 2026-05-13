package releaseproof

import "crypto/sha256"

type manifestData struct {
	manifest Manifest
	ref      string
	digest   [sha256.Size]byte
}

type verificationState struct {
	sourceCommit   string
	state          string
	commitStatus   string
	artifactStatus string
	artifactCounts ArtifactCounts
	artifactIssues []ArtifactIssue
	sourceReason   string
}
