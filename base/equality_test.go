package base

import "testing"

func TestEqResults(t *testing.T) {
	result1 := CheckResult{
		Check: Check{
			Description: "Mismatch between tbl_one и tbl_two",
			Query:       "SELECT * FROM tbl;",
		},
		Problems: []Row{
			Row{"181620", "4", "15"},
			Row{"236695", "2", "3"},
		},
	}
	result2 := CheckResult{
		Check: Check{
			Description: "Other check",
			Query:       "SELECT * FROM tbl;",
		},
		Problems: []Row{
			Row{"181620", "-200", "t1"},
		},
	}
	result3 := CheckResult{
		Check: Check{
			Description: "Mismatch between tbl_one и tbl_two",
			Query:       "SELECT * FROM tbl;",
		},
		Problems: []Row{
			Row{"181620", "4", "15"},
			Row{"236695", "2", "3"},
		},
	}

	if eqResult(result1, result2) {
		t.Errorf("Expected %v to be not equal to %v", result1, result2)
	}
	if !eqResult(result1, result3) {
		t.Errorf("Expected %v to be equal to %v", result1, result3)
	}
}
