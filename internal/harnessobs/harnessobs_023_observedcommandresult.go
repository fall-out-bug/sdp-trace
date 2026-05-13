package harnessobs

import (
	"time"
)

type observedCommandResult struct {
	waitErr error
	end     time.Time
}
