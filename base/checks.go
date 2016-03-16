package base

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Check is a description of check
type Check struct {
	Description string `yaml:"description"`
	Query       string `yaml:"query"`
	Assert      string `yaml:"assert"`
}

// ReadCheck reads check from io.Reader
func ReadCheck(f io.Reader) (*Check, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	c := Check{}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	if c.Description == "" {
		return nil, errors.New("not a valid check, 'description' is missing")
	}
	if c.Query == "" {
		return nil, errors.New("not a valid check, 'query' is missing")
	}
	if c.Assert == "" {
		return nil, errors.New("not a valid check, 'assert' is missing")
	}
	return &c, err
}

// ReadCheckFile reads check from file at filePath.
func ReadCheckFile(filePath string) (*Check, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadCheck(f)
}

// CheckQueryAbsent is a checker function that considers any output row a problem
func CheckQueryAbsent(db *sql.DB, check Check) (*CheckResult, error) {
	var results []Row

	rows, err := db.Query(check.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	rawResult := make([][]byte, len(cols))

	fields := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		fields[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}
	for rows.Next() {
		err = rows.Scan(fields...)
		if err != nil {
			Error.Println(err)
			continue
		}
		results = append(results, NewRow(fields))

	}
	err = rows.Err() // Get any error encountered during iteration
	if err != nil {
		return nil, err
	}
	return &CheckResult{Check: check, Problems: results, Columns: cols}, nil
}

// CheckQueryPresent is a checker function that considers missing output a problem
func CheckQueryPresent(db *sql.DB, check Check) (*CheckResult, error) {
	var results []Row

	rows, err := db.Query(check.Query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		results = append(results, Row{"No results found"})
	}

	return &CheckResult{Check: check, Problems: results}, nil
}

// GetChecks scans filesistem under searchDir and returns list of checks
func GetChecks(searchDir string) ([]*Check, error) {
	var results []*Check

	stat, err := os.Stat(searchDir)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("No a directory: %s", searchDir)
	}

	err = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if ext == ".yaml" || ext == ".yml" {
			check, err := ReadCheckFile(path)
			if err != nil {
				Error.Printf("Failed to read check %s: %v", f.Name(), err)
				return nil
			}
			results = append(results, check)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// RunChecks runs all check, uses dbConnString for db connection
func RunChecks(dbConnString string, checks []*Check, concurrency int) ([]CheckResult, error) {
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return runChecks(db, checks, concurrency)
}

// runChecks runs all checks, uses db object
func runChecks(db *sql.DB, checks []*Check, concurrency int) ([]CheckResult, error) {
	var results []CheckResult
	if concurrency < 1 {
		concurrency = 1
	}

	// channel to get check results
	ch := make(chan *CheckResult)
	defer close(ch)
	// use this channel as a semaphore to limit concurrency
	sem := make(chan bool, concurrency)
	defer close(sem)
	// iterating is less error-prone and helps DRY
	for _, check := range checks {
		// spawn gorutine
		go func(c *Check) {
			var checker func(*sql.DB, Check) (*CheckResult, error)
			// Try to get semaphore. If it is full, we'll block until some other goroutine will end
			sem <- true
			// defer releasing of semaphore
			defer func() { <-sem }()
			switch c.Assert {
			case "absent":
				checker = CheckQueryAbsent
			case "present":
				checker = CheckQueryPresent
			default:
				ch <- FailedCheck(c, fmt.Sprintf("Unknown check assertion %s", c.Assert))
				return
			}
			// perform check
			cr, err := checker(db, *c)
			if err != nil {
				ch <- FailedCheck(c, fmt.Sprintf("Error while running check: %v", err))
				return
			}
			// send result to channel
			ch <- cr
		}(check)
	}

	// get the results
	for _ = range checks {
		cr := <-ch
		results = append(results, *cr)
	}
	// suck all remaining values from sem
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	return results, nil
}
