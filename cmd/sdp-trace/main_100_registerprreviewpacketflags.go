package main

func registerPRReviewPacketFlags(opts *flagSet) {
	// Packet metadata is fully flag-driven so generated review packets can be
	// replayed without hidden process context.
	for _, flag := range prReviewPacketStringFlags {
		opts.setString(flag.name, flag.defaultValue)
	}
}
