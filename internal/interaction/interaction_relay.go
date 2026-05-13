package interaction

import (
	"context"
	"errors"
	"io"
)

func Relay(ctx context.Context, opts RelayOptions, stdin io.Reader, stdout, stderr io.Writer) (Trace, int, error) {
	// Relay keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	opts = normalizeRelay(opts)
	if len(opts.Command) == 0 {
		return Trace{}, 0, errors.New("interaction relay requires forward command after --")
	}
	return relayWithCommand(ctx, opts, stdin, stdout, stderr)
}

func relayWithCommand(ctx context.Context, opts RelayOptions, stdin io.Reader, stdout, stderr io.Writer) (Trace, int, error) {
	// relayWithCommand keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	body, err := readBody(stdin)
	if err != nil {
		return Trace{}, 0, err
	}
	event, err := NewObservedEvent(opts, body, 0)
	if err != nil {
		return Trace{}, 0, err
	}
	trace := NewTrace(opts.TaskID, SourceObservedControlChannel, []Event{event}, opts.Now)
	if err := writeRelayTrace(opts.Out, trace, event, body); err != nil {
		return Trace{}, 0, err
	}
	exitCode, err := runForward(ctx, opts.Command, body, stdout, stderr)
	return trace, exitCode, err
}

func writeRelayTrace(path string, trace Trace, event Event, body []byte) error {
	// writeRelayTrace keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if err := WriteContentBlobs(path, trace, map[string][]byte{event.InteractionID: body}); err != nil {
		return err
	}
	return WriteTrace(path, trace)
}
