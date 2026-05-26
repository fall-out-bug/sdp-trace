package main

import "context"

func runOPAPolicyEval(ctx context.Context, fixtureName string, expectPass bool, label string) (verifierState, string) {
	regoPath, fixturePath, err := opaFixturePaths(fixtureName)
	if err != nil {
		return stateCannotVerify, err.Error()
	}
	pass, err := runOPAEval(ctx, regoPath, fixturePath, "data.sdp_trace.adapter.pass")
	if err != nil {
		return stateFail, err.Error()
	}
	return assertOPAResult(pass, expectPass, label)
}
