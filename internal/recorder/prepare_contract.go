package recorder

import (
	"errors"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// Contract preparation binds the run to a named verifier contract before any
// command evidence is emitted. A missing identifier is a setup failure, not an
// unverifiable run result.

func resolveRecorderContract(contractPath string, useDefault bool) (trace.Contract, error) {
	// A recorder run must bind to a named contract because later replay and
	// provenance checks key evidence to that contract identifier.
	contract, err := resolveContract(contractPath, useDefault)
	if err != nil {
		return trace.Contract{}, err
	}
	if contract.ContractID == "" {
		return trace.Contract{}, errors.New("contract missing identifier")
	}
	return contract, nil
}

func initializeRunWriter(writer *runWriter, options RecorderOptions, contract trace.Contract) error {
	// External contract files are copied into manifest metadata by reference and
	// digest, which lets later verifier code detect drift without rerecording.
	if options.ContractPath != "" {
		writer.manifest.ContractPath = options.ContractPath
		writer.manifest.ContractDigest = trace.SHA256Hex(string(mustMarshalJSON(contract)))
	}
	return writer.writeManifest()
}
