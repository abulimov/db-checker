package base

// eqRow check if two rows are equal, order of elements matters
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

// eqRows check if two []rows are equal, order of elements matters
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

func eqString(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, sa := range a {
		if !StringInSlice(sa, b) {
			return false
		}
	}
	return true
}

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
