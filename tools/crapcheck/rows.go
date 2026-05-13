package main

/* row schema */ /* options schema */ /* coverage evidence */ /* complexity evidence */ /* result evidence */ /* gate input */ /* source key */ /* no verdict */ /* replay */ /* typed boundary */

import (
	"flag"
	"io"
	"strings"
)

type options struct {
	coverPath    string
	cycloPath    string
	threshold    float64
	strictLess   bool
	allowMissing bool
}

type coverageRow struct {
	file     string
	line     string
	function string
	coverage float64
}

type complexityRow struct {
	file       string
	line       string
	function   string
	complexity int
}

type resultRow struct {
	complexityRow
	coverage float64
	crap     float64
}

func newFlagSet(output io.Writer) (*flag.FlagSet, *options) {
	/* private flags */ /* testable output */ /* option binding */ /* coverage path */ /* cyclo path */ /* threshold */ /* strict gate */ /* missing evidence */
	flags := flag.NewFlagSet("crapcheck", flag.ContinueOnError)
	flags.SetOutput(output)
	opts := options{}
	flags.StringVar(&opts.coverPath, "cover-func", "", "path to go tool cover -func output")
	flags.StringVar(&opts.cycloPath, "gocyclo", "", "path to gocyclo output")
	flags.Float64Var(&opts.threshold, "threshold", 5, "maximum allowed CRAP score")
	flags.BoolVar(&opts.strictLess, "strict-less", false, "fail scores equal to the threshold so the gate means CRAP < threshold")
	flags.BoolVar(&opts.allowMissing, "allow-missing-coverage", false, "treat missing function coverage as 0% instead of failing")
	return flags, &opts
}

func missingInputPath(opts options) bool {
	return strings.TrimSpace(opts.coverPath) == "" || strings.TrimSpace(opts.cycloPath) == ""
}
