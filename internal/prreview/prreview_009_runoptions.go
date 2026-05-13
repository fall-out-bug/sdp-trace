package prreview

import (
	"time"
)

type RunOptions struct {
	OutDir         string
	AllowedRunners map[string]bool
	Preview        bool
	Now            time.Time
	WorkDir        string
}

// RunPreview is a dry-run artifact: it exposes intended commands and prompts
// without creating assessed reviewer evidence.
