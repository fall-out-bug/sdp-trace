package trace

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LoadContract returns a parsed contract from path.
func LoadContract(path string) (Contract, error) {
	// Empty path explicitly selects the portable default contract rather than a
	// missing external spec file.
	if path == "" {
		return DefaultContract, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Contract{}, err
	}
	contract, err := parseContract(data)
	if err != nil {
		return Contract{}, err
	}
	return contract.withDefaults(path), nil
}

func parseContract(data []byte) (Contract, error) {
	// Keep parsing logic intentionally light for first milestone; full schema
	// validation is handled in phase-8 schema migration work.
	var contract Contract
	if err := json.Unmarshal(data, &contract); err != nil {
		// JSON parse failure means the contract cannot be trusted as a spec.
		return Contract{}, err
	}
	return contract, nil
}

func (contract Contract) withDefaults(path string) Contract {
	// Contract defaults preserve older lightweight fixtures while still naming
	// the loaded file as the contract ID when no explicit ID is present.
	if contract.ContractID == "" {
		// The path-derived ID is display context, not external authority.
		contract.ContractID = filepath.Base(path)
	}
	if contract.Version == "" {
		// Default version keeps legacy fixtures parseable under the current
		// schema contract.
		contract.Version = SchemaVersion
	}
	if len(contract.RequiredEvents) == 0 {
		// Copy the default slice so callers cannot mutate shared contract state.
		contract.RequiredEvents = append([]string(nil), DefaultContract.RequiredEvents...)
	}
	return contract
}
