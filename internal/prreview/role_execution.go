package prreview

import (
	"time"
)

// runRole orchestrates a single profile role from initial result metadata
// through runner preparation, command execution, completion, and raw capture.
func runRole(packet Packet, role ReviewRole, opts RunOptions, rawDir string) (ReviewerResult, error) {
	result := newReviewerResult(packet, role, opts.Now)
	result.CommandDigest = commandDigest(role.Command)
	baseline, ready, err := prepareRoleRunner(&result, role, opts)
	if err != nil || !ready {
		return result, err
	}
	output, timedOut, err := runRoleCommand(packet, role, opts)
	result = finishRoleRun(result, role, packet, opts, baseline, output, timedOut, err)
	return writeRawResult(result, role, rawDir, output)
}

func finishRoleRun(result ReviewerResult, role ReviewRole, packet Packet, opts RunOptions, baseline *workingTreeBaseline, output []byte, timedOut bool, err error) ReviewerResult {
	result.EndedAt = time.Now().UTC().Format(time.RFC3339)
	return completeRoleResult(result, role, packet, opts.WorkDir, baseline, output, timedOut, err)
}
