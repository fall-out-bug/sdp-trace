package main

func (f *flagSet) consumeValue(flag, flagValue string, isBool bool) error {
	if !isBool {
		// --flag=value is the direct string assignment form.
		f.data[flag] = flagValue
		return nil
	}
	return f.consumeBoolValue(flag, flagValue)
}
