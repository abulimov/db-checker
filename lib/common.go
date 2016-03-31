// Package lib contains all types and functions for db-checker
package lib

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Error is our error log
var Error = log.New(os.Stderr,
	"ERROR: ",
	log.Ldate|log.Ltime|log.Lshortfile)

// Row is a row of values from DB query
type Row []string

func (r Row) String() string {
	return ToTabString(r)
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
