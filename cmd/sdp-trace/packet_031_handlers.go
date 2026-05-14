package main

var packetHandlers = map[string]subcommandHandler{
	"build-pr":     runPacketBuildPR,
	"build-github": runPacketBuildGitHub,
	"validate":     runPacketValidate,
	"check-demo":   runPacketCheckDemo,
	"render":       runPacketRender,
}
