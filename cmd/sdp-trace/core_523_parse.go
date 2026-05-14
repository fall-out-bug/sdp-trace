package main

func (f *flagSet) parse(args []string) error {
	rest := make([]string, 0)
	for i := 0; i < len(args); i++ {
		// The loop index is passed by pointer so string flags can consume their
		// following value without reparsing it as positional input.
		// consumeArg owns index advancement for flags with following values.
		done, err := f.consumeArg(args, &i, &rest)
		if err != nil {
			return err
		}
		if done {
			f.args = rest
			return nil
		}
	}
	f.args = rest
	return nil
}
