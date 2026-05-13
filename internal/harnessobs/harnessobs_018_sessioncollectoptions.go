package harnessobs

import (
	"time"
)

type SessionCollectOptions struct {
	ProfilePath string
	RunDir      string
	Now         time.Time
}
