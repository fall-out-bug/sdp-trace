package main

var interactionRelayRequiredFlags = []requiredCLIFlag{
	{"task-id", "interaction relay requires --task-id"},
	{"event-type", "interaction relay requires --event-type"},
	{"out", "interaction relay requires --out"},
}

var interactionImportTranscriptRequiredFlags = []requiredCLIFlag{
	{"source", "interaction import-transcript requires --source"},
	{"task-id", "interaction import-transcript requires --task-id"},
	{"events-jsonl", "interaction import-transcript requires --events-jsonl"},
	{"out", "interaction import-transcript requires --out"},
}

var interactionSummarizeRequiredFlags = []requiredCLIFlag{
	{"trace", "interaction summarize requires --trace"},
}
