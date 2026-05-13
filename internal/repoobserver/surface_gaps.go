package repoobserver

const ReasonManualStepRequired = "manual_step_required"

func gapsFor(surfaces []Surface) []Gap {
	// Every non-passing surface remains visible as a gap with a concrete reason.
	gaps := make([]Gap, 0)
	for _, s := range surfaces {
		gap, ok := gapForSurface(s)
		if !ok {
			continue
		}
		gaps = append(gaps, gap)
	}
	return gaps
}

func gapForSurface(s Surface) (Gap, bool) {
	// Agent-prompt gaps get a custom explanation because prompt cooperation is
	// not repository setup proof.
	if s.InstallState == StatePass && s.ProofState == StatePass {
		return Gap{}, false
	}
	if agentPromptNotAssessed(s) {
		return agentPromptGap(s), true
	}
	return Gap{SurfaceID: s.SurfaceID, ReasonCode: s.ReasonCode, Detail: gapDetail(s)}, true
}

func agentPromptNotAssessed(s Surface) bool {
	// A prompt-only claim is always a gap because it is not replayable
	// repository evidence.
	return s.InstallState == StateNotAssessed && s.ProofState == StateNotAssessed && s.SurfaceID == SurfaceAgentPrompt
}

func agentPromptGap(s Surface) Gap {
	return Gap{SurfaceID: s.SurfaceID, ReasonCode: s.ReasonCode, Detail: "agent prompt cooperation is not repository setup proof"}
}

func gapDetail(s Surface) string {
	if s.NextAction != "" {
		// Prefer concrete remediation over terse reason codes in human summaries.
		return s.NextAction
	}
	return s.ReasonCode
}
