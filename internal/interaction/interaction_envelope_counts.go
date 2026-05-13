package interaction

func envelopeReferenceCounts(envelope Envelope) envelopeRefCounts {
	// envelopeReferenceCounts keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	counts := primaryEnvelopeReferenceCounts(envelope)
	addExecutionEnvelopeReferenceCounts(&counts, envelope)
	return counts
}

func primaryEnvelopeReferenceCounts(envelope Envelope) envelopeRefCounts {
	// primaryEnvelopeReferenceCounts keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	run := len(envelope.RunRefs)
	source := len(envelope.SourceRefs)
	task := len(envelope.TaskRefs)
	promise := len(envelope.PromiseRefs)

	return envelopeRefCounts{run: run, source: source, task: task, promise: promise}
}

func addExecutionEnvelopeReferenceCounts(counts *envelopeRefCounts, envelope Envelope) {
	// addExecutionEnvelopeReferenceCounts keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	counts.interaction = len(envelope.InteractionRefs)
	counts.operation = len(envelope.OperationRefs)
	counts.llm = len(envelope.LLMRefs)
	counts.tool = len(envelope.ToolRefs)

	counts.mutation = len(envelope.MutationRefs)
	counts.stage = len(envelope.StageRefs)
	counts.friction = len(envelope.FrictionRefs)
}
