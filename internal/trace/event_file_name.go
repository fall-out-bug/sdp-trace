package trace

import "fmt"

func EventFileName(sequence int, eventType EventType) string {
	// Event filenames mirror sequence order while the hash stays payload-bound.
	return fmt.Sprintf("%06d-%s.json", sequence, eventType)
}
