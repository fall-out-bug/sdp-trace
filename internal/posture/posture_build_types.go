package posture

import (
	"time"
)

type buildInput struct {
	selection  SelectionManifest
	activeKeys []string
	cutoff     time.Time
	hasCutoff  bool
	handoff    map[string]string
}
