package releaseproof

import "time"

const (
	SchemaVersion = "0.1.0"
	TrustScope    = "source_bound_local_release"

	StatePass         = "pass"
	StateFail         = "fail"
	StateCannotVerify = "cannot_verify"
	StatusMatched     = "matched"
	StatusMismatch    = "mismatch"
	StatusMissing     = "missing"
	StatusNotAssessed = "not_assessed"
)

type verificationInput struct {
	manifestData
	verificationState
	verificationTime time.Time
}

func Evaluate(repoRoot, manifestPath string, now time.Time) (Verification, error) {
	// Evaluation has three trust phases: load bounded manifest evidence,
	// compare it to the source commit, then render a conservative verdict.
	manifestData, err := loadManifest(repoRoot, manifestPath)
	if err != nil {
		return Verification{}, err
	}
	state := evaluateManifestState(repoRoot, manifestData.manifest)
	return buildVerification(verificationInput{
		manifestData:      manifestData,
		verificationState: state,
		verificationTime:  now,
	}), nil
}
