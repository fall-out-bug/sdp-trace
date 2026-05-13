package harnessobs

import (
	"time"
)

type observationContext struct {
	outDir       string
	sourcePath   string
	sourceDigest string
	now          time.Time
	profile      Profile
	events       []Event
}
