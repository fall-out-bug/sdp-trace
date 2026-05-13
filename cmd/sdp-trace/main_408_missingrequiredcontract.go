package main

func missingRequiredContract(opts *flagSet) bool {
	return opts.stringValue("contract") == "" && !opts.boolValue("use-default-contract")
}
