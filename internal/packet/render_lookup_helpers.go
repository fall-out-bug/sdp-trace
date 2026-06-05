package packet

import "strings"

func rowByID(rows []Row, id string) Row {
	// rowByID keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, row := range rows {
		if row.ID == id {

			return row
		}
	}
	return Row{ID: id, State: StateCannotVerify, Summary: "row missing", Reason: "row missing"}
}

func requiredRowIndex(id string) int {
	// requiredRowIndex keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for i, required := range RequiredRows {
		if id == required {

			return i
		}
	}
	return len(RequiredRows)
}

func resolverFromList(resolvers []ResolverEntry, ref string) string {
	// resolverFromList keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	for _, resolver := range resolvers {
		if resolver.Ref == ref {

			return resolver.Resolver
		}
	}
	return ""
}

func md(value string) string {
	// md keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "|", "\\|")
	if strings.TrimSpace(value) == "" {
		return "none"
	}
	return value
}
