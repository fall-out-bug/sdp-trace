package main

type githubActionsArtifactPayload struct {
	Artifacts []githubActionsArtifact `json:"artifacts"`
}
