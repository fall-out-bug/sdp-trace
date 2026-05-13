package harnessobs

import (
	"time"
)

type ObserveOptions struct {
	ProfilePath string
	SourcePath  string
	OutDir      string
	Now         time.Time
}
