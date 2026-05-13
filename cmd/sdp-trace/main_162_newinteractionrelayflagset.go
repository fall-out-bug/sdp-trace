package main

func newInteractionRelayFlagSet() *flagSet {
	opts := &flagSet{name: "interaction relay"}
	// Relay defaults encode a human-to-agent corrective-feedback event; callers
	// override them only when the trace source is more specific.
	for _, flag := range interactionRelayStringFlags {
		opts.setString(flag.name, flag.defaultValue)
	}
	return opts
}
