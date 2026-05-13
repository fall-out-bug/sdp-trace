package packet

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
