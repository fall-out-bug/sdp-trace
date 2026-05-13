package packet

import (
	"bytes"

	"fmt"
)

func renderNoResidualGaps(out *bytes.Buffer) {

	fmt.Fprintf(out, "No residual gaps recorded beyond row states.\n\n")
}
