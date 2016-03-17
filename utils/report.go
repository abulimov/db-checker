// Package utils contains all functions to read/write/generate reports for db-checker
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/tabwriter"

	"github.com/abulimov/db-checker/base"
)

// ReportProblems counts and pretty print problems
func ReportProblems(results []base.CheckResult) (int, string) {
	w := new(tabwriter.Writer)
	var buffer *bytes.Buffer
	count := 0
	prettyCount := 0
	report := ""
	prettyNumbers := false
	for _, cr := range results {
		buffer = new(bytes.Buffer)
		w.Init(buffer, 1, 1, 0, ' ', 0)
		prettyCount = 0
		prettyNumbers = false
		if cr.HasProblems() {
			report += "\n* " + cr.Check.Description + "\n"
			if len(cr.Columns) != 0 {
				prettyNumbers = true
				fmt.Fprintf(w, "N. \t¦ %s\n", base.ToTabString(cr.Columns))
			}
			for _, p := range cr.Problems {
				count++
				prettyCount++
				if prettyNumbers {
					fmt.Fprintf(w, "%d. \t¦ %s\n", prettyCount, p)
				} else {
					fmt.Fprintf(w, "%s\n", p)
				}
			}
		}
		w.Flush()
		report += buffer.String()
	}
	if count == 0 {
		return count, "No problems found"
	}
	return count, report
}

// WriteReportFile writes report to file at filePath.
func WriteReportFile(results []base.CheckResult, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		base.Error.Println(err)
		return err
	}
	defer f.Close()
	return WriteReport(results, f)

}

// WriteReport writes report in CSV format to some io.Writer.
func WriteReport(results []base.CheckResult, f io.Writer) error {
	data, err := json.Marshal(results)
	if err != nil {
		return err
	}
	io.WriteString(f, string(data))
	return nil
}

// ReadReportFile reads report from file at filePath.
func ReadReportFile(filePath string) ([]base.CheckResult, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadReport(f)
}

// ReadReport reads previous report in CSV format from some io.Reader.
func ReadReport(f io.Reader) ([]base.CheckResult, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var report []base.CheckResult
	err = json.Unmarshal(b, &report)
	if err != nil {
		return nil, err
	}
	return report, nil
}

// DiffResults return diff (in form of []base.CheckResult) between two []base.CheckResult slices
func DiffResults(first, second []base.CheckResult) []base.CheckResult {
	var add []base.CheckResult

	for _, s := range second {
		pos := base.FindCheckInCheckResults(s.Check, first)
		if pos == -1 {
			add = append(add, s)
		} else {
			old := first[pos]
			diff := base.DiffRows(old.Problems, s.Problems)
			if len(diff) > 0 {
				add = append(add, base.CheckResult{
					Check:    s.Check,
					Columns:  s.Columns,
					Problems: diff,
				})
			}
		}
	}
	return add
}
