package main

import (
	"context"
	"fmt"
	"time"
)

// runInTotoWrap tests in-toto-run presence.
func runInTotoWrap() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := runExternalTool(ctx, "in-toto-run", "--version"); err != nil {
		return stateCannotVerify, fmt.Sprintf("in-toto-run version check failed: %v", err)
	}
	return stateCannotVerify, "in-toto-run present; run manual wrap per docs/oss-replacement-compatibility.md"
}

// runCosignLocalSign tests cosign presence.
func runCosignLocalSign() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := runExternalTool(ctx, "cosign", "version"); err != nil {
		return stateCannotVerify, fmt.Sprintf("cosign version check failed: %v", err)
	}
	return stateCannotVerify, "cosign present; run manual sign/verify per docs/oss-replacement-compatibility.md"
}

// runSLSANegative tests slsa-verifier presence.
func runSLSANegative() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if _, err := runExternalTool(ctx, "slsa-verifier", "version"); err != nil {
		return stateCannotVerify, fmt.Sprintf("slsa-verifier version check failed: %v", err)
	}
	return stateCannotVerify, "slsa-verifier present; run manual negative test per docs/oss-replacement-compatibility.md"
}
