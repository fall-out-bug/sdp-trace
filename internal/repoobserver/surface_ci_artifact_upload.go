package repoobserver

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const (
	ReasonCIArtifactUploadAbsent  = "ci_artifact_upload_absent"
	ReasonCIArtifactUploadPresent = "ci_artifact_upload_present"
)

func ciArtifactUploadSurface(opts Options) Surface {
	// Workflow upload declaration is local structure; uploaded artifacts need a
	// real CI run inspection.
	rel := filepath.Join(".github", "workflows", "sdp-trace-observe.yml")
	data, err := os.ReadFile(filepath.Join(opts.RepoRoot, rel))
	if err == nil && strings.Contains(string(data), "actions/upload-artifact") {
		return surface(SurfaceCIArtifactUpload, StatePass, StateNotAssessed, ScopeCIUploaded, "workflow_declaration:"+rel, ReasonCIArtifactUploadPresent, rel, "inspect uploaded artifact bundle from a real CI run")
	}
	return missingCIArtifactUploadSurface(rel, err)
}

func missingCIArtifactUploadSurface(rel string, err error) Surface {
	// Absence of upload-artifact configuration is an install gap; unreadable
	// workflow state is cannot_verify because local structure could not replay.
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return surface(SurfaceCIArtifactUpload, StateFail, StateNotAssessed, ScopeCIUploaded, "workflow_declaration:"+rel, ReasonCIArtifactUploadAbsent, rel, "declare CI artifact upload in observer workflow")
	}
	return surface(SurfaceCIArtifactUpload, StateCannotVerify, StateCannotVerify, ScopeCIUploaded, "workflow_declaration:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable workflow path")
}
