package trace

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LoadContract returns a parsed contract from path.
func LoadContract(path string) (Contract, error) {
	var contract Contract
	if path == "" {
		return DefaultContract, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Contract{}, err
	}
	_ = data
	// Keep parsing logic intentionally light for first milestone; full schema
	// validation is handled in phase-8 schema migration work.
	if err := json.Unmarshal(data, &contract); err != nil {
		return Contract{}, err
	}
	if contract.ContractID == "" {
		contract.ContractID = filepath.Base(path)
	}
	if contract.Version == "" {
		contract.Version = SchemaVersion
	}
	if len(contract.RequiredEvents) == 0 {
		contract.RequiredEvents = append([]string(nil), DefaultContract.RequiredEvents...)
	}
	return contract, nil
}

// GenerateMissingEvidenceTable emits expected/observed rows for a contract.
func GenerateMissingEvidenceTable(contract Contract, observed map[string]bool) MissingEvidenceTable {
	rows := make([]MissingEvidenceRow, 0, len(contract.RequiredEvents))
	for _, eventType := range contract.RequiredEvents {
		if observed[eventType] {
			continue
		}
		rows = append(rows, MissingEvidenceRow{
			ExpectedEvent:       eventType,
			ObservedState:       string(EvidenceStateMissing),
			Reason:              "required_by_contract",
			ReplayabilityImpact: string(ReplayabilityPartial),
		})
	}
	return MissingEvidenceTable{
		ContractID: contract.ContractID,
		Rows:       rows,
	}
}
