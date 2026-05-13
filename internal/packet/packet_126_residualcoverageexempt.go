package packet

func residualCoverageExempt(row Row) bool {

	return row.ID == "PC-RESIDUAL-GAPS" || row.State == StatePass
}
