package main

import "flag"

func registerFlags(fs *flag.FlagSet) (*bool, *bool, *string) {
	asJSON := fs.Bool("json", false, "emit JSON output")
	list := fs.Bool("list", false, "list available probes")
	probe := fs.String("probe", "", "run a single probe by name")
	return asJSON, list, probe
}
