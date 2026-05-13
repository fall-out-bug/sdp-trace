package trace

// EnsureDefaults populates event defaults used during hashing and writing.
func (event Event) EnsureDefaults() Event {
	// Defaults are applied before payload synchronization so hashes include the
	// same schema, hash, and canonicalization fields that will be persisted.
	event = event.withDefaultIdentityFields()
	event = event.withDefaultCanonicalFields()
	if synced, err := event.syncPayloadRepresentation(); err == nil {
		event = synced
	}
	return event
}

func (event Event) withDefaultIdentityFields() Event {
	// Empty hash algorithm and previous hash fields are legacy-compatible
	// defaults, not evidence claims from the caller.
	if event.SchemaVersion == "" {
		event.SchemaVersion = SchemaVersion
	}
	if event.HashAlgorithm == "" {
		event.HashAlgorithm = HashAlgSHA256
	}
	if event.PrevEventHash == "" {
		event.PrevEventHash = NullEventHash
	}
	return event
}

func (event Event) withDefaultCanonicalFields() Event {
	// Canonicalization defaults bind new events to the in-repo canonical JSON
	// contract.
	if event.Canonicalization.Algorithm == "" {
		event.Canonicalization.Algorithm = CanonicalSchemaAlgo
	}
	if event.Canonicalization.Version == "" {
		event.Canonicalization.Version = CanonicalAlgoVersion
	}
	return event
}
