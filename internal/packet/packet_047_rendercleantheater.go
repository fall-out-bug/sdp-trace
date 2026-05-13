package packet

import (
	"bytes"

	"fmt"
)

func renderCleanTheater(out *bytes.Buffer, theater Row) {

	fmt.Fprintf(out, "| none | %s | none | %s | PC-THEATER row | %s |\n\n", theater.State, md(theater.Summary), md(theater.Reason))
}
