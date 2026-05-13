package main

func newCheckpointCreateFlagSet() *flagSet {
	opts := &flagSet{name: "checkpoint create"}
	for _, flag := range checkpointCreateStringFlags {
		// Registration order is fixed so help/tests observe the same command
		// contract while defaults stay beside their flag names.
		opts.setString(flag.name, flag.defaultValue)
	}
	return opts
}
