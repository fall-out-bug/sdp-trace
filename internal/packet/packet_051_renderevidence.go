package packet

import (
	"bytes"

	"fmt"
)

func renderEvidence(out *bytes.Buffer, manifest BundleManifest) {
	// renderEvidence keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	fmt.Fprintf(out, "## Evidence Bundle\n\n")
	fmt.Fprintf(out, "Manifest: `%s`\n\n", md(manifest.BundleID))
	fmt.Fprintf(out, "| ref | source class | retained form | redaction status | resolver |\n| --- | --- | --- | --- | --- |\n")
	for _, entry := range manifest.Entries {

		resolver := entry.Resolver
		if resolver == "" {

			resolver = resolverFromList(manifest.Resolvers, entry.Ref)
		}
		fmt.Fprintf(out, "| %s | %s | %s | %s | %s |\n", md(entry.Ref), md(entry.SourceClass), md(entry.RetainedForm), md(entry.RedactionStatus), md(resolver))
	}
	fmt.Fprintln(out)
}
