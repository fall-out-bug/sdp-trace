package witness

const (
	// ReasonUnsupported records a profile that the current witness verifier does
	// not support.
	ReasonUnsupported = "witness_unsupported_profile"
	// ReasonUnsafeOutput records output that would expose unsafe data classes.
	ReasonUnsafeOutput = "witness_unsafe_output_candidate"
	// ReasonPrivateKeyInput records rejected private-key input material.
	ReasonPrivateKeyInput = "witness_private_key_input_rejected"
	// ReasonMalformedInput records malformed witness input data.
	ReasonMalformedInput = "witness_malformed_input"
)
