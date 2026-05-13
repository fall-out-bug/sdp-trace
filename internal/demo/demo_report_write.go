package demo

import (
	"errors"
	"strings"
)

func WriteReport(target, outDir, contractPath string) (ReportArtifacts, error) {
	// WriteReport keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if strings.TrimSpace(outDir) == "" {

		return ReportArtifacts{}, errors.New("report requires --out <dir>")
	}
	rows, contract, err := verifiedRowsForContract(target, contractPath)
	if err != nil {
		return ReportArtifacts{}, err
	}
	artifacts := BuildReport(rows, contract)
	if err := persistReportArtifacts(outDir, artifacts); err != nil {
		return ReportArtifacts{}, err
	}
	return artifacts, nil
}
