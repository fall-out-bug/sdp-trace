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
