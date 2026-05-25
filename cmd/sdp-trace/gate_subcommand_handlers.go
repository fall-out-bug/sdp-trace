package main

var gateSubcommandHandlers = map[string]subcommandHandler{
	"explain": runGateExplain,
	"preview": runGatePreview,
}
