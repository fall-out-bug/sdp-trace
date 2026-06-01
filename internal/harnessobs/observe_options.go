package harnessobs

func requireObserveOptions(opts ObserveOptions) error {
	if err := requireNonBlank(opts.ProfilePath, "harness observe requires --profile"); err != nil {
		return err
	}
	if err := requireNonBlank(opts.SourcePath, "harness observe requires --source"); err != nil {
		return err
	}
	if err := requireNonBlank(opts.OutDir, "harness observe requires --out"); err != nil {
		return err
	}
	return nil
}
