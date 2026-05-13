package repoobserver

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	ReasonCIWorkflowAbsent  = "ci_workflow_absent"
	ReasonCIWorkflowPresent = "ci_workflow_present"
)

func ciWorkflowSurface(opts Options) Surface {
	// A checked-in workflow is not proof that CI has executed it.
	rel := filepath.Join(".github", "workflows", "sdp-trace-observe.yml")
	_, err := os.ReadFile(filepath.Join(opts.RepoRoot, rel))
	if err == nil {
		reason := ReasonCIWorkflowPresent
		proof := StateNotAssessed
		return surface(SurfaceCIWorkflow, StatePass, proof, ScopeLocalStructural, "filesystem:"+rel, reason, rel, "observe a CI run artifact before treating workflow as proof")
	}
	return missingCIWorkflowSurface(rel, err)
}

func missingCIWorkflowSurface(rel string, err error) Surface {
	if errors.Is(err, os.ErrNotExist) {
		// Missing workflow is an install gap; unreadable workflow paths below are
		// cannot_verify because the local filesystem state could not be replayed.
		return surface(SurfaceCIWorkflow, StateFail, StateNotAssessed, ScopeLocalStructural, "filesystem:"+rel, ReasonCIWorkflowAbsent, rel, "install GitHub Actions observer workflow")
	}
	return surface(SurfaceCIWorkflow, StateCannotVerify, StateCannotVerify, ScopeLocalStructural, "filesystem:"+rel, ReasonUnsafeOutputRefused, rel, "fix unreadable workflow path")
}
