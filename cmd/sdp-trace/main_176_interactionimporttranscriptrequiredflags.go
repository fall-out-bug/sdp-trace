package main

var interactionImportTranscriptRequiredFlags = []requiredCLIFlag{
	{"task-id", "interaction import-transcript requires --task-id"},
	{"events-jsonl", "interaction import-transcript requires --events-jsonl"},
	{"out", "interaction import-transcript requires --out"},
}
