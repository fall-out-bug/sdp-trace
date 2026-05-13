package main

func newCheckpointVerifyFlagSet() *flagSet {
	opts := &flagSet{name: "checkpoint verify"}
	// Verification names the run, signed checkpoint, and optional trust policy
	// as separate inputs so each can fail independently.
	for _, flag := range checkpointVerifyStringFlags {
		opts.setString(flag, "")
	}
	return opts
}
