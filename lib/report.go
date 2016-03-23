package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

// ReportProblems counts and pretty print problems
func ReportProblems(results []CheckResult) (int, string) {
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
				fmt.Fprintf(w, "N. \t¦ %s\n", ToTabString(cr.Columns))
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
func WriteReportFile(results []CheckResult, filePath string) error {
	tmpfile, err := ioutil.TempFile("", "db-checker")
	if err != nil {
		return err
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if err = WriteReport(results, tmpfile); err != nil {
		return err
	}
	if err = tmpfile.Close(); err != nil {
		return err
	}
	if err = os.Rename(tmpfile.Name(), filePath); err != nil {
		return err
	}
	return nil

}

// WriteReport writes report in CSV format to some io.Writer.
func WriteReport(results []CheckResult, f io.Writer) error {
	data, err := json.Marshal(results)
	if err != nil {
		return err
	}
	io.WriteString(f, string(data))
	return nil
}

// ReadReportFile reads report from file at filePath.
func ReadReportFile(filePath string) ([]CheckResult, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadReport(f)
}

// ReadReport reads previous report in CSV format from some io.Reader.
func ReadReport(f io.Reader) ([]CheckResult, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var report []CheckResult
	err = json.Unmarshal(b, &report)
	if err != nil {
		return nil, err
	}
	return report, nil
}
