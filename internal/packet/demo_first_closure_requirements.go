package packet

func (c *demoFirstPacketChecker) requireVerificationOrReviewAssessed() {
	if rowAssessed(c.rows["PC-VERIFICATION"]) || rowAssessed(c.rows["PC-REVIEW"]) {
		return
	}

	c.add("demo first-packet gate requires PC-VERIFICATION or PC-REVIEW to be pass, partial, or fail")
}

func rowAssessed(row Row) bool {
	return row.State == StatePass || row.State == StatePartial || row.State == StateFail
}
