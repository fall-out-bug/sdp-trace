package interaction

func SummarizeEnvelope(envelope Envelope) Summary {

	return envelopeSummary(envelope, envelopeReferenceCounts(envelope))
}

func envelopeSummary(envelope Envelope, counts envelopeRefCounts) Summary {
	// envelopeSummary keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	return Summary{
		SchemaVersion: SchemaVersion,

		TaskID:     envelope.TaskID,
		EnvelopeID: envelope.EnvelopeID,

		AssessmentState: envelope.AssessmentState,

		RunRefCount:         counts.run,
		SourceRefCount:      counts.source,
		TaskRefCount:        counts.task,
		PromiseRefCount:     counts.promise,
		InteractionRefCount: counts.interaction,
		OperationRefCount:   counts.operation,
		LLMRefCount:         counts.llm,
		ToolRefCount:        counts.tool,
		MutationRefCount:    counts.mutation,
		StageRefCount:       counts.stage,
		FrictionRefCount:    counts.friction,

		NotAssessed:  envelope.NotAssessed,
		CannotVerify: envelope.CannotVerify,
	}
}
