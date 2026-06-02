package main

var checkpointCreateStringFlags = []struct {
	name         string
	defaultValue string
}{
	{"run", ""},
	{"out", ""},
	{"private-key", ""},
	{"signer-id", ""},
	{"id", "checkpoint-001"},
}

var checkpointCreateRequiredFlags = []requiredCLIFlag{
	{"run", "checkpoint create requires --run"},
	{"out", "checkpoint create requires --out"},
	{"private-key", "checkpoint create requires --private-key"},
	{"signer-id", "checkpoint create requires --signer-id"},
}
