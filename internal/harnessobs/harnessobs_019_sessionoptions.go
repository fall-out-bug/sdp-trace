package harnessobs

import (
	"time"
)

type SessionOptions struct {
	ProfilePath string
	OutDir      string
	Command     []string
	Now         time.Time
}
