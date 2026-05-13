package trace

import (
	"fmt"
	"regexp"
	"strconv"
)

// EventSeqFromFilename returns the numeric sequence prefix from `<NNNNNN>-<event>.json`.
func EventSeqFromFilename(name string) (int, error) {
	// The filename sequence is a sorting hint only; event-chain validation still
	// checks the sequence field inside the event.
	re := regexp.MustCompile(`^(\d{6})-`)
	matches := re.FindStringSubmatch(name)
	if len(matches) != 2 {
		return 0, fmt.Errorf("invalid event filename: %s", name)
	}
	return strconv.Atoi(matches[1])
}
