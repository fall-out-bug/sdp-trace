package repoobserver

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	ReasonCIArtifactBundleNotObserved = "ci_artifact_bundle_not_observed"
	ReasonCIArtifactBundleObserved    = "ci_artifact_bundle_observed"
)

func ciArtifactBundleSurface(opts Options) Surface {
	// Local CI artifact directories are only structural unless downloaded and
	// bound to CI artifact storage.
	rel := filepath.Join(".sdp-trace", "ci")
	path := filepath.Join(opts.RepoRoot, rel)
	entries, err := os.ReadDir(path)
	if err == nil && len(entries) > 0 {
		return surface(SurfaceCIArtifactBundleObservation, StatePass, StateNotAssessed, ScopeCIUploaded, "filesystem:"+rel, ReasonCIArtifactBundleObserved, rel, "treat as local structural only unless downloaded from CI artifact storage")
	}
	return missingCIArtifactBundleSurface(rel, err)
}

func missingCIArtifactBundleSurface(rel string, err error) Surface {
	// Missing local CI artifacts is not a failure by itself; it means CI evidence
	// has not been inspected for this profile.
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return surface(SurfaceCIArtifactBundleObservation, StateNotAssessed, StateNotAssessed, ScopeCIUploaded, "filesystem:"+rel, ReasonCIArtifactBundleNotObserved, rel, "run CI and inspect uploaded artifact bundle")
	}
	return surface(SurfaceCIArtifactBundleObservation, StateCannotVerify, StateCannotVerify, ScopeCIUploaded, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable CI artifact path")
}
