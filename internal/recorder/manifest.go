package recorder

import (
	"path/filepath"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func initialRunManifest(contract trace.Contract, task string) trace.RunManifest {
	// The contract digest is computed from the resolved contract object because
	// the recorder may use either a file-backed or built-in contract.
	manifest := baseRunManifest(contract.ContractID, task)
	manifest.ContractDigest = trace.SHA256Hex(string(mustMarshalJSON(contract)))
	return manifest
}

func baseRunManifest(contractID, task string) trace.RunManifest {
	// Mutable fields start empty or unknown and are filled only after matching
	// evidence exists in the event chain or source snapshot.
	return trace.RunManifest{
		SchemaVersion:   trace.SchemaVersion,
		RunID:           randomHex(16),
		RecorderVersion: trace.RecorderVersion,
		CreatedAt:       time.Now().UTC().Format(time.RFC3339Nano),

		Task:       task,
		ContractID: contractID,

		SourceSnapshot: "",
		SourceState:    "",

		EventCount:     0,
		ClosureState:   trace.ClosureStateUnknown,
		ContractPath:   "",
		ContractDigest: "",
	}
}

func (w *runWriter) writeManifest() error {
	// Chain-head fields are written only after at least one event exists, which
	// avoids presenting an empty run as if it had replayable closure evidence.
	w.manifest.EventCount = w.sequence
	if w.sequence > 0 {
		w.manifest.EventChainHead = w.lastHash
		w.manifest.FinalChainHead = w.lastHash
	}
	return writeIndentedJSON(filepath.Join(w.runDir, "run.json"), w.manifest)
}
