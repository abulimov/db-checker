package base

// eqRow check if two Rows are equal, order of elements matters
func eqRow(first, second Row) bool {
	if len(first) != len(second) {
		return false
	}
	for i, f := range first {
		if second[i] != f {
			return false
		}
	}
	return true
}

// eqRows check if two []Rows are equal, order of elements matters
func eqRows(first, second []Row) bool {
	if len(first) != len(second) {
		return false
	}
	for i, f := range first {
		if !eqRow(second[i], f) {
			return false
		}
	}
	return true
}

// eqResult checks if tow CheckResults are equal
func eqResult(a, b CheckResult) bool {
	if a.Check != b.Check {
		return false
	}
	if !eqRow(a.Columns, b.Columns) {
		return false
	}
	if !eqRows(a.Problems, b.Problems) {
		return false
	}
	return true
}
