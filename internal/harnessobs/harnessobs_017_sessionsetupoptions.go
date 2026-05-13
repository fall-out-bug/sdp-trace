package harnessobs

import (
	"time"
)

type SessionSetupOptions struct {
	ProfilePath string
	OutDir      string
	Command     string
	Now         time.Time
}
