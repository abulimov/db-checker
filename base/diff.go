package base

import (
	"bytes"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffRows calculates diff between slices of Row
func DiffRows(first, second []Row) []Row {
	var add []Row
	m := make(map[string]int)

	for _, y := range first {
		s := y.String()
		if _, ok := m[s]; ok {
			m[s]++
		} else {
			m[s] = 1
		}
	}

	for _, x := range second {
		if m[x.String()] > 0 {
			m[x.String()]--
			continue
		}
		add = append(add, x)
	}
	return add
}

// DiffResults return diff between two results slices
func DiffResults(first, second []CheckResult) []CheckResult {
	var add []CheckResult

	for _, s := range second {
		pos := FindCheckInCheckResults(s.Check, first)
		if pos == -1 {
			add = append(add, s)
		} else {
			old := first[pos]
			diff := DiffRows(old.Problems, s.Problems)
			if len(diff) > 0 {
				add = append(add, CheckResult{
					Check:    s.Check,
					Columns:  s.Columns,
					Problems: diff,
				})
			}
		}
	}
	return add
}

// DiffPretty returns pretty diff between two strings
func DiffPretty(a, b string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(a, b, true)
	var buff bytes.Buffer
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buff.WriteString("{+>")
			buff.WriteString(diff.Text)
			buff.WriteString("-}")
		case diffmatchpatch.DiffDelete:
			buff.WriteString("{-")
			buff.WriteString(diff.Text)
			buff.WriteString("}")
		case diffmatchpatch.DiffEqual:
			buff.WriteString("")
			buff.WriteString(diff.Text)
			buff.WriteString("")
		}
	}
	return buff.String()
}
