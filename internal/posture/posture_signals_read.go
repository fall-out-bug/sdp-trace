package posture

import (
	"fmt"
	"strings"
)

func readSignals(path string) (map[string]PostureSignal, error) {
	// readSignals keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if strings.TrimSpace(path) == "" {

		return map[string]PostureSignal{}, nil
	}
	manifest, err := readSignalManifest(path)
	if err != nil {
		return nil, err
	}
	return validatedSignals(manifest.Signals)
}

func readSignalManifest(path string) (SignalManifest, error) {
	// readSignalManifest keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	manifest, err := readJSONFile[SignalManifest](path)
	if err != nil {
		return manifest, err
	}

	if manifest.SchemaVersion != SignalManifestSchemaVersion {
		return manifest, fmt.Errorf("unsupported signal manifest schema")
	}
	return manifest, nil
}

func validatedSignals(signals []PostureSignal) (map[string]PostureSignal, error) {
	// validatedSignals keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	out := map[string]PostureSignal{}
	for _, signal := range signals {
		if unsafeSignal(signal) {
			return nil, fmt.Errorf("unsafe signal")
		}
		out[signal.RowRef] = signal
	}
	return out, nil
}

// unsafeSignal checks whether a posture signal crosses the output safety boundary
// before incorporation into metric evidence. Unsafe keywords block signal incorporation.
func unsafeSignal(signal PostureSignal) bool {
	return unsafeOutput(signal.RowRef + signal.WitnessScope + signal.ObserverState + signal.OverrideMarker + signal.LateAttachMarker + signal.ContractChangeMarker)
}
