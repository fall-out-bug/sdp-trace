package main

var packetHandlers = map[string]subcommandHandler{
	"build-pr":     runPacketBuildPR,
	"build-github": runPacketBuildGitHub,
	"validate":     runPacketValidate,
	"check-demo":   runPacketCheckDemo,
	"render":       runPacketRender,
}

var packetBuildPRRequiredFlags = []requiredCLIFlag{
	{"out", "packet build-pr requires --out"},
}

var packetBuildGitHubRequiredFlags = []requiredCLIFlag{
	{"github-input", "packet build-github requires --github-input"},
	{"out", "packet build-github requires --out"},
}

var packetValidateRequiredFlags = []requiredCLIFlag{
	{"bundle", "packet validate requires --bundle"},
}

var packetCheckDemoRequiredFlags = []requiredCLIFlag{
	{"bundle", "packet check-demo requires --bundle"},
}

var packetRenderRequiredFlags = []requiredCLIFlag{
	{"bundle", "packet render requires --bundle"},
	{"out", "packet render requires --out"},
}
