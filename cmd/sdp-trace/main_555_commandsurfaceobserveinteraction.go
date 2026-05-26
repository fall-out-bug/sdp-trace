package main

var commandSurfaceInteraction = commandSurfaceCmd{
	Name:        "interaction",
	Description: "Interaction recording and summarization.",
	Subcommands: []string{"relay", "import-transcript", "summarize"},
	Variations: []string{
		"sdp-trace interaction relay --task-id \u003csafe-id\u003e --event-type \u003ctype\u003e --out \u003cfile\u003e -- \u003cforward-command...\u003e",
		"sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id \u003csafe-id\u003e --events-jsonl \u003cfile\u003e --out \u003cfile\u003e",
		"sdp-trace interaction summarize --trace \u003cfile\u003e [--out \u003cfile\u003e]",
	},
	State: "partial",
}
