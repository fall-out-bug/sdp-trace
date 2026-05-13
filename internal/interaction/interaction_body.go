package interaction

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

func readBody(stdin io.Reader) ([]byte, error) {
	// Empty stdin is rejected because a relay event without content would have no
	// observed interaction payload to hash.
	// The limit reader detects oversize bodies without buffering unbounded
	// interaction content.
	var buf bytes.Buffer
	limited := io.LimitReader(stdin, MaxBodyBytes+1)
	if _, err := buf.ReadFrom(limited); err != nil {
		return nil, err
	}
	body := buf.Bytes()
	if len(body) == 0 {
		return nil, errors.New("interaction relay requires stdin content")
	}
	if len(body) > MaxBodyBytes {
		return nil, fmt.Errorf("interaction content exceeds %d bytes", MaxBodyBytes)
	}
	return body, nil
}
