package main

func publicGitHubServerHost(serverHost string) bool {
	return serverHost == "" || serverHost == "github.com"
}
