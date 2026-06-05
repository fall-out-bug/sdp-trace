package prreview

import "strings"

func unavailablePacketFields(opts PacketOptions) []UnavailableField {
	// unavailablePacketFields keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	fields := []UnavailableField{}
	if strings.TrimSpace(opts.MetadataPath) == "" {

		fields = append(fields, UnavailableField{Field: "metadata_ref", State: StateNotAssessed, Reason: "metadata_input_not_provided"})
	}
	if len(opts.ContextPaths) == 0 {

		fields = append(fields, UnavailableField{Field: "context_refs", State: StateNotAssessed, Reason: "context_inputs_not_provided"})
	}
	if len(opts.VerificationPaths) == 0 {

		fields = append(fields, UnavailableField{Field: "verification_refs", State: StateNotAssessed, Reason: "verification_inputs_not_provided"})
	}
	return fields
}
