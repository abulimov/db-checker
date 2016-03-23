package lib

import "testing"

func TestDiffRows(t *testing.T) {
	first := []Row{
		Row{"181620", "4", "15"},
		Row{"181620", "4", "15"},
		Row{"236695", "2", "3"},
		Row{"381630", "43", "153"},
	}
	second := []Row{
		Row{"181620", "4", "15"},
		Row{"181620", "4", "15"},
		Row{"236695", "2", "3"},
		Row{"581650", "53", "1"},
	}

	expectedAdd := []Row{
		Row{"581650", "53", "1"},
	}

	add := DiffRows(first, second)
	if len(add) != len(expectedAdd) {
		t.Fatalf("Expected diff to have len %d, got %d", len(expectedAdd), len(add))
	}
}

func TestDiffResults(t *testing.T) {
	t1 := "2015-08-10 11:42:50.641621+03"
	first := []CheckResult{
		{
			Check: Check{
				Description: "Mismatch between tbl_one и tbl_two",
				Query:       "SELECT * from tbl",
			},
			Problems: []Row{
				{"181620", "4", "15"},
				{"236695", "2", "3"},
			},
		},
		{
			Check: Check{
				Description: "Other check",
				Query:       "SELECT * from tbl LIMIT 1",
			},
			Problems: []Row{
				{"181620", "-200", t1},
			},
		},
	}
	second := []CheckResult{
		{
			Check: Check{
				Description: "Mismatch between tbl_one и tbl_two",
				Query:       "SELECT * from tbl",
			},
			Problems: []Row{
				{"181620", "4", "15"},
				{"998899", "20", "30"},
			},
		},
		{
			Check: Check{
				Description: "Other check",
				Query:       "SELECT * from tbl LIMIT 1",
			},
			Problems: []Row{
				{"181620", "-200", t1},
			},
		},
	}
	expectedAdd := []CheckResult{
		{
			Check: Check{
				Description: "Mismatch between tbl_one и tbl_two",
				Query:       "SELECT * from tbl",
			},
			Problems: []Row{
				{"998899", "20", "30"},
			},
		},
	}

	add := DiffResults(first, second)

	if len(add) != len(expectedAdd) {
		t.Errorf("Got add len %v not equal to expected %v", len(add), len(expectedAdd))
	}

	for _, p := range add {
		if !ResultInSlice(p, expectedAdd) {
			t.Errorf("Result %v not found in expected add", p)
		}
	}
}
