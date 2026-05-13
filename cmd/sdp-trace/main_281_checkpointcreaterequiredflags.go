package main

var checkpointCreateRequiredFlags = []requiredCLIFlag{
	{"run", "checkpoint create requires --run"},
	{"out", "checkpoint create requires --out"},
	{"private-key", "checkpoint create requires --private-key"},
	{"signer-id", "checkpoint create requires --signer-id"},
}
