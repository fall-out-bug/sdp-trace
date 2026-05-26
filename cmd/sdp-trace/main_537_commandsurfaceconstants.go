package main

const commandSurfaceSchemaVersion = "sdp-trace-command-surface-v1"

func knownAssessmentProfiles() []profileMeta {
	return commandSurfaceProfiles
}

func knownWitnessKinds() []string {
	return commandSurfaceWitnessKinds
}

func knownStates() []stateMeta {
	return commandSurfaceStates
}
