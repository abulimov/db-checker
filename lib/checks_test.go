package lib

import (
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestReadCheckBad(t *testing.T) {
	data := `
	{
		"description": "Some description",
		"something": "test",
	}
	`
	_, err := ReadCheck(strings.NewReader(data))
	if err == nil {
		t.Fatal("Expected to fail to read bad check")
	}
}

func TestReadCheckOk(t *testing.T) {
	data := `
description: Some description
query: SELECT * FROM some_table
assert: present
`
	gotCheck, err := ReadCheck(strings.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to read check: %v", err)
	}
	expectedCheck := Check{
		Description: "Some description",
		Query:       "SELECT * FROM some_table",
		Assert:      "present",
	}

	if *gotCheck != expectedCheck {
		t.Errorf("Got check %v not equal to expected %v", gotCheck, expectedCheck)
	}
}

func TestQueryPresentOk(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	check := Check{
		Description: "some_table non-empty",
		Query:       "SELECT id, some_col FROM some_table",
		Assert:      "present",
	}

	columns := []string{"id", "some_col"}
	// match query it with regexp
	mock.ExpectQuery(`SELECT id, some_col FROM some_table`).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1, OK"))

	result, err := CheckQueryPresent(db, check)
	if err != nil {
		t.Fatalf("Expected no error, but got %s instead", err)
	}
	gotResult := result
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	if result.HasProblems() {
		t.Error("Expected result to have no problems")
	}

	expectedLen := 0
	gotLen := len(gotResult.Problems)
	if gotLen != expectedLen {
		t.Errorf("Expected len of problems %v, got %v", expectedLen, gotLen)
	}
}

func TestQueryPresentBad(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	check := Check{
		Description: "some_table non-empty",
		Query:       "SELECT id, some_col FROM some_table",
		Assert:      "present",
	}

	columns := []string{"id", "some_col"}
	// match query it with regexp
	mock.ExpectQuery(`SELECT id, some_col FROM some_table`).
		WillReturnRows(sqlmock.NewRows(columns))

	result, err := CheckQueryPresent(db, check)
	if err != nil {
		t.Fatalf("Expected no error, but got %s instead", err)
	}

	expectedResult := CheckResult{
		Check: check,
		Problems: []Row{
			{"No results found"},
		},
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	if !eqResult(*result, expectedResult) {
		t.Errorf("Expected result to have expected problems %v, got %v", *result, expectedResult)
	}
}

func TestQueryAbsentOk(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	check := Check{
		Description: "some_table empty",
		Query:       "SELECT id, some_col FROM some_table",
		Assert:      "absent",
	}
	// match query it with regexp
	columns := []string{"id", "some_col"}
	mock.ExpectQuery(`SELECT id, some_col FROM some_table`).
		WillReturnRows(sqlmock.NewRows(columns))

	result, err := CheckQueryAbsent(db, check)
	if err != nil {
		t.Fatalf("Expected no error, but got %s instead", err)
	}
	gotResult := result
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	if result.HasProblems() {
		t.Error("Expected result to have no problems")
	}

	expectedLen := 0
	gotLen := len(gotResult.Problems)
	if gotLen != expectedLen {
		t.Errorf("Expected len of problems %v, got %v", expectedLen, gotLen)
	}
}

func TestQueryAbsentBad(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	check := Check{
		Description: "some_table empty",
		Query:       "SELECT id, some_col FROM some_table",
		Assert:      "absent",
	}
	columns := []string{"id", "some_col"}
	// match query it with regexp
	mock.ExpectQuery(`SELECT id, some_col FROM some_table`).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1, OK"))

	result, err := CheckQueryAbsent(db, check)
	if err != nil {
		t.Fatalf("Expected no error, but got %s instead", err)
	}

	expectedResult := CheckResult{
		Check:   check,
		Columns: Row{"id", "some_col"},
		Problems: []Row{
			{"1", "OK"},
		},
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	if !eqResult(*result, expectedResult) {
		t.Errorf("Expected result to have expected problems %v, got %v", *result, expectedResult)
	}
}

func TestGetChecks(t *testing.T) {
	badPath := "./no_such_dir/"

	_, err := GetChecks(badPath)
	if err == nil {
		t.Fatalf("Expected to fail to read checks from bad path %v", badPath)
	}

	testPath := "./test_data/"

	checks, err := GetChecks(testPath)
	if err != nil {
		t.Fatalf("Failed to get checks from path %s: %v", testPath, err)
	}

	expectedCheck := Check{
		Description: "some_table empty",
		Query:       "SELECT id, some_col FROM some_table",
		Assert:      "absent",
	}
	expectedLen := 1

	if len(checks) != expectedLen {
		t.Errorf("Expected to find %d check, found %d", expectedLen, len(checks))
	}

	if *checks[0] != expectedCheck {
		t.Errorf("Expected check to be equal to %v, found %v", expectedCheck, *checks[0])
	}
}

func TestRunChecks(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.MatchExpectationsInOrder(false)

	checks := []*Check{
		{
			Description: "some_table empty",
			Query:       "SELECT id, some_col FROM some_table",
			Assert:      "absent",
		},
		{
			Description: "other_table non-empty",
			Query:       "SELECT id, other_col FROM other_table",
			Assert:      "present",
		},
	}
	columns1 := []string{"id", "some_col"}
	columns2 := []string{"id", "other_col"}
	// match query it with regexp
	mock.ExpectQuery(`SELECT id, some_col FROM some_table`).
		WillReturnRows(sqlmock.NewRows(columns1).FromCSVString("1, OK"))
	mock.ExpectQuery(`SELECT id, other_col FROM other_table`).
		WillReturnRows(sqlmock.NewRows(columns2))

	result, err := runChecks(db, checks, 1)
	if err != nil {
		t.Fatalf("Expected no error, but got %s instead", err)
	}

	expectedResult := []CheckResult{
		{
			Check:   *checks[0],
			Columns: Row{"id", "some_col"},
			Problems: []Row{
				{"1", "OK"},
			},
		},
		{
			Check: *checks[1],
			Problems: []Row{
				{"No results found"},
			},
		},
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expections: %s", err)
	}

	if len(result) != len(expectedResult) {
		t.Fatalf("Expected %d results, got %d", len(expectedResult), len(result))
	}

	for _, r := range expectedResult {
		if !ResultInSlice(r, result) {
			t.Errorf("Expected to find %v in results %v", r, result)
		}
	}
}

func TestCheckQueryBool(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	check := Check{
		Description: "some_table count more than 5",
		Query:       "SELECT (count(*) > 5) as status FROM some_table",
		Assert:      "true",
	}

	columns := []string{"status"}
	// match query it with regexp
	mock.ExpectQuery(`SELECT.+`).
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("t"))

	result, err := CheckQueryBool(db, check, true)
	if err != nil {
		t.Fatalf("Expected no error, but got %s instead", err)
	}
	gotResult := result
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	if result.HasProblems() {
		t.Error("Expected result to have no problems")
	}

	expectedLen := 0
	gotLen := len(gotResult.Problems)
	if gotLen != expectedLen {
		t.Errorf("Expected len of problems %v, got %v", expectedLen, gotLen)
	}
}

func TestCheckQueryBoolEmpty(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	check := Check{
		Description: "some_table count more than 5",
		Query:       "SELECT (count(*) > 5) as status FROM some_table",
		Assert:      "true",
	}
	// match query it with regexp
	columns := []string{"status"}
	mock.ExpectQuery(`SELECT.+`).
		WillReturnRows(sqlmock.NewRows(columns))

	result, err := CheckQueryBool(db, check, false)
	if err != nil {
		t.Fatalf("Expected no error, but got %s instead", err)
	}
	expectedResult := CheckResult{
		Check: check,
		Problems: []Row{
			{"No rows for boolean check"},
		},
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	if !eqResult(*result, expectedResult) {
		t.Errorf("Expected result to have expected problems %v, got %v", *result, expectedResult)
	}
}
