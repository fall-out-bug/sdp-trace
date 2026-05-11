package prreview

import (
	"errors"
	"regexp"
)

var (
	errPromptEvidenceCannotVerify = errors.New("prompt_evidence_cannot_verify")
	repoIDPattern                 = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,62}[a-z0-9]$`)
	changeRefPattern              = regexp.MustCompile(`^(pr|mr|change)-[A-Za-z0-9._-]{1,64}$`)
	sha40Pattern                  = regexp.MustCompile(`^[0-9a-f]{40}$`)
)
