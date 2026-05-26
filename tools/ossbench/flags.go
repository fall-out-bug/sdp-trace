package main

import (
	"flag"
	"io"
)

func parseFlagsAndArgs(args []string, stderr io.Writer) (runConfig, error) {
	fs := flag.NewFlagSet("ossbench", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = usageFunc(fs, stderr)
	var cfg runConfig
	registerFlags(fs, &cfg)
	if err := fs.Parse(args); err != nil {
		return runConfig{}, err
	}
	cfg.args = fs.Args()
	return cfg, nil
}
