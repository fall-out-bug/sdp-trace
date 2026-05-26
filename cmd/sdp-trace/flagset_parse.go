package main

func (f *flagSet) parse(args []string) error {
	rest := make([]string, 0)
	for i := 0; i < len(args); i++ {
		// Each iteration either stores positional input or delegates a complete
		// flag token, preserving argv order for command-specific validation.
		done, err := f.consumeArg(args, &i, &rest)
		if err != nil {
			return err
		}
		if done {
			// A stop signal means "--" moved the remaining payload into rest.
			f.args = rest
			return nil
		}
	}
	f.args = rest
	return nil
}
