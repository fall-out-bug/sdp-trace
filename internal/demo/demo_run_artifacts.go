package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func applyRunArtifact(row *RunRow, artifact trace.RunArtifact, contract trace.Contract) {
	// applyRunArtifact keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	row.RunID = artifact.Manifest.RunID
	row.ClosureState = artifact.Manifest.ClosureState
	commandStarted, commandFinished := commandEvents(artifact.Events)
	row.OverrideRequests = overrideRequestsFromEvents(artifact.Events, contract)

	row.Command = payloadString(commandStarted, "command")
	row.WrapperName = payloadString(commandStarted, "wrapper_name")
	if exitCode, ok := payloadInt(commandFinished, "exit_code"); ok {

		row.ExitCode = &exitCode
	}
	row.StdoutDigest = payloadString(commandFinished, "stdout_digest")
	row.StderrDigest = payloadString(commandFinished, "stderr_digest")

	row.Kind, row.KindReason = classify(artifact.Events, contract.RequiredEvidence)
}
