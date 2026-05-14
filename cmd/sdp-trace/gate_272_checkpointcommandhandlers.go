package main

var checkpointCommandHandlers = map[string]subcommandHandler{
	"create": runCheckpointCreate,
	"verify": runCheckpointVerify,
}
