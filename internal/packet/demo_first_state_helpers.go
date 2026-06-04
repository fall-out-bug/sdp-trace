package packet

func rowAssessed(row Row) bool {
	return row.State == StatePass || row.State == StatePartial || row.State == StateFail
}
