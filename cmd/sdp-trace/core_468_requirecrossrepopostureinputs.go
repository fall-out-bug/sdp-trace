package main

import (
	"fmt"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

func requireCrossRepoPostureInputs(opts *flagSet) error {
	if strings.TrimSpace(opts.stringValue("profile")) != posture.ProfileID {
		// Keep the CLI bound to the only supported cross-repo posture contract.
		return fmt.Errorf("export cross-repo-posture requires --profile cross-repo-evidence-posture-v1")
	}
	if strings.TrimSpace(opts.stringValue("selection")) == "" {
		// The selection artifact is the auditable repository set for posture.
		return fmt.Errorf("export cross-repo-posture requires --selection")
	}
	return nil
}
