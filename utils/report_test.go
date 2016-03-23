package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/abulimov/db-checker/base"
)

var t1 = "2015-08-10 11:42:50.641621+03"
var exampleContent = `[{"check":{"Description":"Mismatch between tbl_one and tbl_two","Query":"SELECT * FROM tbl;","Assert":""},"problems":[["181620","4","15"],["236695","2","30"]],"columns":["ID","F","S"]},{"check":{"Description":"Other check","Query":"SELECT * FROM tbl;","Assert":""},"problems":[["181620","-200","2015-08-10 11:42:50.641621+03"]],"columns":["user_id","balance","date"]},{"check":{"Description":"Another check","Query":"SELECT * FROM tbl;","Assert":""},"problems":[["Warner Bros. Entertainment, Inc.","Interview with the Vampire: The Vampire Chronicles","vampire","2015-08-10 11:42:50.641621+03"],["Sony Pictures","Repentance","some-slug","2015-08-10 11:42:50.641621+03"]],"columns":["rightsholder","title","slug","date"]}]`
var exampleCheckResults = []base.CheckResult{
	{
		Check: base.Check{
			Description: "Mismatch between tbl_one and tbl_two",
			Query:       "SELECT * FROM tbl;",
		},
		Columns: base.Row{"ID", "F", "S"},
		Problems: []base.Row{
			{"181620", "4", "15"},
			{"236695", "2", "30"},
		},
	},
	{
		Check: base.Check{
			Description: "Other check",
			Query:       "SELECT * FROM tbl;",
		},
		Columns: base.Row{"user_id", "balance", "date"},
		Problems: []base.Row{
			{"181620", "-200", t1},
		},
	},
	{

		Check: base.Check{
			Description: "Another check",
			Query:       "SELECT * FROM tbl;",
		},
		Columns: base.Row{"rightsholder", "title", "slug", "date"},
		Problems: []base.Row{
			{"Warner Bros. Entertainment, Inc.", "Interview with the Vampire: The Vampire Chronicles", "vampire", t1},
			{"Sony Pictures", "Repentance", "some-slug", t1},
		},
	},
}

func TestReportProblems(t *testing.T) {
	results := exampleCheckResults
	gotCount, gotReport := ReportProblems(results)
	expectedCount := 5
	if gotCount != expectedCount {
		t.Errorf("Expected count of problems %v, got %v", expectedCount, gotCount)
	}

	expectedReport := `
* Mismatch between tbl_one and tbl_two
N. ¦ ID     ¦ F ¦ S
1. ¦ 181620 ¦ 4 ¦ 15
2. ¦ 236695 ¦ 2 ¦ 30

* Other check
N. ¦ user_id ¦ balance ¦ date
1. ¦ 181620  ¦ -200    ¦ ` + t1 + `

* Another check
N. ¦ rightsholder                     ¦ title                                              ¦ slug      ¦ date
1. ¦ Warner Bros. Entertainment, Inc. ¦ Interview with the Vampire: The Vampire Chronicles ¦ vampire   ¦ ` + t1 + `
2. ¦ Sony Pictures                    ¦ Repentance                                         ¦ some-slug ¦ ` + t1 + `
`

	if gotReport != expectedReport {
		t.Errorf(
			"Diff between actual and expected reports:\n'%v'",
			base.DiffPretty(gotReport, expectedReport),
		)
	}
}

func TestWriteReport(t *testing.T) {
	results := exampleCheckResults
	expectedContent := exampleContent
	var gotBytes bytes.Buffer
	err := WriteReport(results, &gotBytes)
	if err != nil {
		t.Fatalf("Got error %v on writing report", err)
	}
	gotContent := gotBytes.String()
	if expectedContent != gotContent {
		t.Errorf(
			"Diff between actual and expected reports:\n'%v'",
			base.DiffPretty(gotContent, expectedContent),
		)
	}
}

func TestReadReport(t *testing.T) {
	expectedResults := exampleCheckResults
	content := exampleContent
	gotResults, err := ReadReport(strings.NewReader(content))
	if err != nil {
		t.Fatalf("Got error %v on reading report", err)
	}
	gotLen := len(gotResults)
	expectedLen := len(expectedResults)
	if expectedLen != gotLen {
		t.Errorf("Got results len %v not equal to expected %v", gotLen, expectedLen)
	}
	for _, p := range gotResults {
		if !base.ResultInSlice(p, expectedResults) {
			t.Errorf("Result %v not found in expected results", p)
		}
	}
}
