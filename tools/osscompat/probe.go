package main

// verifierState is the canonical probe result state.
type verifierState string

const (
	statePass         verifierState = "pass"
	stateFail         verifierState = "fail"
	stateCannotVerify verifierState = "cannot_verify"
	stateNotAssessed  verifierState = "not_assessed"
)

// probeResult is the output of a single compatibility probe.
type probeResult struct {
	Name   string        `json:"name"`
	State  verifierState `json:"state"`
	Reason string        `json:"reason,omitempty"`
}

// probe is a single compatibility probe that can be run.
type probe struct {
	Name        string
	NeedsTool   string
	Run         func() (verifierState, string)
	Description string
}

// registry holds all defined probes.
var registry = []probe{
	{
		Name:        "jsonschema-fixtures",
		NeedsTool:   "check-jsonschema",
		Description: "Verify check-jsonschema is present and can validate fixtures (run manually per docs)",
		Run:         runJSONSchemaFixtures,
	},
	{
		// Keep this probe pointed at the live manifest schema. The richer
		// flight-recorder-run schema is a separate profile artifact.
		// The drift-era name remains an alias so old reproduction commands do
		// not silently lose coverage after the semantic split.
		Name:        "jsonschema-wrap-manifest",
		NeedsTool:   "check-jsonschema",
		Description: "Validate live wrap manifest schema",
		Run:         runJSONSchemaWrapManifest,
	},
	{
		Name:        "opa-policy",
		NeedsTool:   "opa",
		Description: "Evaluate adapter.rego against the positive test fixture",
		Run:         runOPAPolicy,
	},
	{
		Name:        "opa-negative",
		NeedsTool:   "opa",
		Description: "Evaluate adapter.rego against the combined negative test fixture",
		Run:         runOPANegativeFixture,
	},
	{
		Name:        "opa-negative-traceid",
		NeedsTool:   "opa",
		Description: "Evaluate adapter.rego against the negative trace_id fixture",
		Run:         runOPANegativeTraceID,
	},
	{
		Name:        "opa-negative-provenance",
		NeedsTool:   "opa",
		Description: "Evaluate adapter.rego against the negative provenance fixture",
		Run:         runOPANegativeProvenance,
	},
	{
		Name:        "cue-import",
		NeedsTool:   "cue",
		Description: "Verify cue can import JSON Schema to stdout without mutating working tree",
		Run:         runCUEImport,
	},
	{
		Name:        "intoto-wrap",
		NeedsTool:   "in-toto-run",
		Description: "Verify in-toto-run is present and responds to version query",
		Run:         runInTotoWrap,
	},
	{
		Name:        "cosign-local-sign",
		NeedsTool:   "cosign",
		Description: "Verify cosign is present and responds to version query",
		Run:         runCosignLocalSign,
	},
	{
		Name:        "slsa-negative",
		NeedsTool:   "slsa-verifier",
		Description: "Verify slsa-verifier is present and responds to version query",
		Run:         runSLSANegative,
	},
}

var legacyProbeNames = map[string]string{
	"jsonschema-wrap-drift": "jsonschema-wrap-manifest",
}
