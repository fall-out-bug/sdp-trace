package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func witnessHasProtectedTrust(witnessSummary demo.WitnessSummary) bool {
	return witnessSummary.Kind == "github-actions" && witnessSummary.Status == demo.GatePass && witnessSummary.TrustScope == "ci_witnessed"
}
