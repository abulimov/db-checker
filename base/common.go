package base

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// PGTime is format used by Postgres to represent date time
const PGTime = "2006-01-02 15:04:05.999999-07"

// Error is our error log
var Error = log.New(os.Stderr,
	"ERROR: ",
	log.Ldate|log.Ltime|log.Lshortfile)

// Row is a row of values from DB query
type Row []string

func (r Row) String() string {
	return ToTabString(r)
}

// StringInSlice checks if string is in slice
func StringInSlice(needle string, list []string) bool {
	for _, b := range list {
		if b == needle {
			return true
		}
	}
	return false
}

// ToCSVString returns keys and values in CSV format
func ToCSVString(keys, values []string) string {
	var list []string
	for pos, key := range keys {
		list = append(list, fmt.Sprintf("%s: %s", key, values[pos]))
	}
	return strings.Join(list, ", ")
}

// ToTabString returns keys and values in CSV format
func ToTabString(values []string) string {
	var list []string
	for _, val := range values {
		list = append(list, strings.TrimSpace(val))
	}
	return strings.Join(list, " \t¦ ")
}

// CheckResult is a result of performed checks
type CheckResult struct {
	Check    Check `json:"check"`
	Problems []Row `json:"problems"`
	Columns  Row   `json:"columns"`
}

// HasProblems indicates that CheckResult has problems
func (c CheckResult) HasProblems() bool {
	return len(c.Problems) > 0
}

// FailedCheck provides easy way to create failed CheckResult
func FailedCheck(c *Check, message string) *CheckResult {
	return &CheckResult{
		Check:    *c,
		Problems: []Row{{message}},
	}
}

// String represents CheckResult as string
func (c CheckResult) String() string {
	result := fmt.Sprintf("Check: %v\n", c.Check)
	result += fmt.Sprintf("Columns: %v\nProblems:\n", c.Columns.String())
	for _, p := range c.Problems {
		result += p.String() + "\n"
	}
	return result
}

// ResultInSlice checks if we have CheckResult in slice of CheckResults
func ResultInSlice(needle CheckResult, list []CheckResult) bool {
	for _, b := range list {
		if eqResult(needle, b) {
			return true
		}
	}
	return false
}

// FindCheckInCheckResults returns position of CheckResult with given Check in []CheckResult
func FindCheckInCheckResults(needle Check, list []CheckResult) int {
	for pos, b := range list {
		if needle == b.Check {
			return pos
		}
	}
	return -1
}

// NewProblem is a problem description constructor
func NewProblem(columns []string, fields []interface{}) string {
	var result []string
	for pos, t := range fields {
		switch t := t.(type) {
		default:
			Error.Printf("unexpected type %T\n", t) // %T prints whatever type t has
		case *[]byte:
			result = append(result, fmt.Sprintf("%s: %s", columns[pos], string(*t)))
		}
	}
	return strings.Join(result, ", ")
}

// NewRow is a row constructor
func NewRow(fields []interface{}) Row {
	var result Row
	for _, t := range fields {
		switch t := t.(type) {
		default:
			Error.Printf("unexpected type %T\n", t) // %T prints whatever type t has
		case *[]byte:
			result = append(result, string(*t))
		}
	}
	return result
}
