package main

var assessSubcommandHandlers = map[string]subcommandHandler{
	"preview": runAssessPreview,
	"explain": runAssessExplain,
}
